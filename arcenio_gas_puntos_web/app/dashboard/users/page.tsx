"use client"

import { useEffect, useState, type FormEvent } from "react"
import { apiGet, apiPost, apiPut, apiPatch } from "@/lib/api"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Button } from "@/components/ui/button"
import { Label } from "@/components/ui/label"
import {
    Dialog,
    DialogContent,
    DialogHeader,
    DialogTitle,
    DialogTrigger,
} from "@/components/ui/dialog"
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from "@/components/ui/select"
import { UserPlus, Pencil, Power, Clock } from "lucide-react"
import { toast } from "sonner"

interface User {
    id: string
    nombres: string
    apellidos: string
    cedula: string
    telefono: string
    direccion: string
    username: string
    role: string
    role_name: string
    is_active: boolean
}

interface Transaction {
    id: string
    transaction_type: string
    points: number
    gallons_amount: number
    processed_by_name?: string
    created_at: string
}

function formatCedula(raw: string): string {
    const d = raw.replace(/\D/g, "").substring(0, 11)
    let result = ""
    for (let i = 0; i < d.length; i++) {
        if (i === 3 || i === 10) result += "-"
        result += d[i]
    }
    return result
}

function telefonoMask(value: string): string {
    const d = value.replace(/\D/g, "").substring(0, 10)
    let result = ""
    if (d.length > 0) result += "(" + d.substring(0, 3)
    if (d.length >= 4) result += ") " + d.substring(3, 6)
    if (d.length >= 7) result += "-" + d.substring(6, 10)
    return result
}

function formatTelefono(raw: string): string {
    if (!raw) return ""
    const d = raw.replace(/\D/g, "")
    if (d.length !== 10) return raw
    return `(${d.substring(0, 3)}) ${d.substring(3, 6)}-${d.substring(6)}`
}

function isAdmin(user: User): boolean {
    const rn = (user.role_name || user.role || "").toLowerCase()
    return rn === "admin" || rn === "superadmin"
}

export default function UsersPage() {
    const [users, setUsers] = useState<User[]>([])
    const [loading, setLoading] = useState(true)
    const [dialogOpen, setDialogOpen] = useState(false)
    const [editing, setEditing] = useState<User | null>(null)
    const [error, setError] = useState("")

    // Form fields
    const [nombres, setNombres] = useState("")
    const [apellidos, setApellidos] = useState("")
    const [cedula, setCedula] = useState("")
    const [telefono, setTelefono] = useState("")
    const [direccion, setDireccion] = useState("")
    const [username, setUsername] = useState("")
    const [password, setPassword] = useState("")
    const [roleId, setRoleId] = useState("3") // "2" = admin, "3" = bombero
    const [formLoading, setFormLoading] = useState(false)

    async function fetchUsers() {
        try {
            const [bomberos, admins] = await Promise.all([
                apiGet<User[]>("/users/bomberos"),
                apiGet<User[]>("/users/admins"),
            ])
            const allBomberos = (bomberos || []).map((u) => ({ ...u, role: u.role || "bombero" }))
            const allAdmins = (admins || []).map((u) => ({ ...u, role: u.role || "admin" }))
            setUsers([...allAdmins, ...allBomberos])
        } catch {
            // silently fail
        } finally {
            setLoading(false)
        }
    }

    useEffect(() => { fetchUsers() }, [])

    function resetForm() {
        setNombres(""); setApellidos(""); setCedula(""); setTelefono("")
        setDireccion(""); setUsername(""); setPassword(""); setRoleId("3")
        setEditing(null); setError("")
    }

    function openEdit(user: User) {
        setEditing(user)
        setNombres(user.nombres)
        setApellidos(user.apellidos)
        setCedula(formatCedula(user.cedula))
        setTelefono(user.telefono ? telefonoMask(user.telefono) : "")
        setDireccion(user.direccion)
        setUsername(user.username)
        setRoleId(isAdmin(user) ? "2" : "3")
        setPassword("")
        setError("")
        setDialogOpen(true)
    }

    function getEndpoint(r: string) {
        return r.toLowerCase() === "admin" || r.toLowerCase() === "superadmin" ? "/users/admins" : "/users/bomberos"
    }

    async function handleSubmit(e: FormEvent) {
        e.preventDefault()
        setFormLoading(true)
        setError("")

        const currentRole = editing ? (isAdmin(editing) ? "admin" : "bombero") : (roleId === "2" ? "admin" : "bombero")

        const body: Record<string, unknown> = {
            nombres: nombres.trim(),
            apellidos: apellidos.trim(),
            cedula: cedula.replace(/\D/g, ""),
            telefono: telefono.replace(/\D/g, ""),
            direccion: direccion.trim(),
            username: username.trim(),
            role_id: parseInt(roleId),
        }
        if (password) body.password = password

        try {
            if (editing) {
                // Update only sends the fields UpdateUserRequest expects
                const updateBody: Record<string, string> = {
                    nombres: nombres.trim(),
                    apellidos: apellidos.trim(),
                    telefono: telefono.replace(/\D/g, ""),
                    direccion: direccion.trim(),
                }
                if (password) updateBody.password = password
                await apiPut(`${getEndpoint(editing.role_name || editing.role)}/${editing.id}`, updateBody)
            } else {
                if (!password) { setError("La clave es requerida"); setFormLoading(false); return }
                body.password = password
                await apiPost(getEndpoint(currentRole), body)
            }
            setDialogOpen(false)
            resetForm()
            setLoading(true)
            await fetchUsers()
            toast.success(editing ? "Usuario actualizado" : "Usuario creado exitosamente")
        } catch (err) {
            const msg = err instanceof Error ? err.message : "Error"
            setError(msg)
            toast.error(msg)
        } finally {
            setFormLoading(false)
        }
    }

    async function toggleActive(user: User) {
        try {
            await apiPatch(`${getEndpoint(user.role_name || user.role)}/${user.id}/active`, { is_active: !user.is_active })
            await fetchUsers()
            toast.success(user.is_active ? "Usuario desactivado" : "Usuario activado")
        } catch {
            toast.error("Error al cambiar estado")
        }
    }

    // --- History modal ---
    const [historyOpen, setHistoryOpen] = useState(false)
    const [historyUser, setHistoryUser] = useState<User | null>(null)
    const [history, setHistory] = useState<Transaction[]>([])
    const [historyLoading, setHistoryLoading] = useState(false)

    async function openHistory(user: User) {
        setHistoryUser(user)
        setHistoryOpen(true)
        setHistoryLoading(true)
        try {
            const txs = await apiGet<Transaction[]>(`/reports/transactions?user_id=${user.id}`)
            setHistory(txs || [])
        } catch {
            setHistory([])
        } finally {
            setHistoryLoading(false)
        }
    }

    return (
        <>
            <div className="space-y-6">
                <div className="flex items-center justify-between">
                    <h2 className="text-2xl font-bold">Personal</h2>
                    <Dialog open={dialogOpen} onOpenChange={(open) => { setDialogOpen(open); if (!open) resetForm() }}>
                        <DialogTrigger asChild>
                            <Button>
                                <UserPlus className="mr-2 h-4 w-4" />
                                Nuevo Usuario
                            </Button>
                        </DialogTrigger>
                        <DialogContent className="sm:max-w-lg">
                            <DialogHeader>
                                <DialogTitle>{editing ? "Editar Usuario" : "Registrar Usuario"}</DialogTitle>
                            </DialogHeader>
                            <form onSubmit={handleSubmit} className="space-y-4 mt-2">
                                <div className="grid grid-cols-2 gap-4">
                                    <div>
                                        <Label>Nombres *</Label>
                                        <Input value={nombres} onChange={(e) => setNombres(e.target.value)} required autoFocus />
                                    </div>
                                    <div>
                                        <Label>Apellidos *</Label>
                                        <Input value={apellidos} onChange={(e) => setApellidos(e.target.value)} required />
                                    </div>
                                </div>
                                <div className="grid grid-cols-2 gap-4">
                                    <div>
                                        <Label>Cédula *</Label>
                                        <Input placeholder="###-#######-#" value={cedula} onChange={(e) => setCedula(formatCedula(e.target.value))} required />
                                    </div>
                                    <div>
                                        <Label>Rol *</Label>
                                        <Select value={roleId} onValueChange={setRoleId} disabled={!!editing}>
                                            <SelectTrigger>
                                                <SelectValue />
                                            </SelectTrigger>
                                            <SelectContent>
                                                <SelectItem value="2">Administrador</SelectItem>
                                                <SelectItem value="3">Bombero</SelectItem>
                                            </SelectContent>
                                        </Select>
                                    </div>
                                </div>
                                <div className="grid grid-cols-2 gap-4">
                                    <div>
                                        <Label>Teléfono</Label>
                                        <Input placeholder="(###) ###-####" value={telefono} onChange={(e) => setTelefono(telefonoMask(e.target.value))} />
                                    </div>
                                    <div>
                                        <Label>Dirección</Label>
                                        <Input value={direccion} onChange={(e) => setDireccion(e.target.value)} />
                                    </div>
                                </div>
                                <div className="grid grid-cols-2 gap-4">
                                    <div>
                                        <Label>Usuario *</Label>
                                        <Input value={username} onChange={(e) => setUsername(e.target.value)} required />
                                    </div>
                                    <div>
                                        <Label>Clave (PIN) {editing ? "(vacío = no cambiar)" : "*"}</Label>
                                        <Input
                                            type="password"
                                            inputMode="numeric"
                                            pattern="[0-9]*"
                                            value={password}
                                            onChange={(e) => setPassword(e.target.value.replace(/\D/g, ""))}
                                            required={!editing}
                                        />
                                    </div>
                                </div>
                                {error && <p className="text-destructive text-sm">{error}</p>}
                                <Button type="submit" disabled={formLoading} className="w-full">
                                    {formLoading ? "Guardando..." : editing ? "Guardar Cambios" : "Crear Usuario"}
                                </Button>
                            </form>
                        </DialogContent>
                    </Dialog>
                </div>

                {/* Table */}
                <Card>
                    <CardContent className="pt-6">
                        {loading ? (
                            <p className="text-muted-foreground text-center py-8">Cargando...</p>
                        ) : users.length === 0 ? (
                            <p className="text-muted-foreground text-center py-8">No hay usuarios registrados</p>
                        ) : (
                            <div className="overflow-x-auto">
                                <table className="w-full text-sm">
                                    <thead>
                                        <tr className="border-b">
                                            <th className="text-left py-2 px-3 font-medium text-muted-foreground">Nombre</th>
                                            <th className="text-left py-2 px-3 font-medium text-muted-foreground">Cédula</th>
                                            <th className="text-left py-2 px-3 font-medium text-muted-foreground hidden md:table-cell">Teléfono</th>
                                            <th className="text-left py-2 px-3 font-medium text-muted-foreground">Usuario</th>
                                            <th className="text-left py-2 px-3 font-medium text-muted-foreground">Rol</th>
                                            <th className="text-left py-2 px-3 font-medium text-muted-foreground">Estado</th>
                                            <th className="text-right py-2 px-3 font-medium text-muted-foreground">Acciones</th>
                                        </tr>
                                    </thead>
                                    <tbody>
                                        {users.map((u) => (
                                            <tr key={u.id} className="border-b last:border-0 hover:bg-muted/50">
                                                <td className="py-2.5 px-3 font-medium">{u.nombres} {u.apellidos}</td>
                                                <td className="py-2.5 px-3 whitespace-nowrap">{formatCedula(u.cedula)}</td>
                                                <td className="py-2.5 px-3 hidden md:table-cell text-muted-foreground whitespace-nowrap">
                                                    {u.telefono ? formatTelefono(u.telefono) : "-"}
                                                </td>
                                                <td className="py-2.5 px-3">{u.username}</td>
                                                <td className="py-2.5 px-3">
                                                    <span className={`inline-flex items-center px-3 py-1 rounded-full text-xs font-medium border ${isAdmin(u) ? "border-purple-200 bg-purple-50 text-purple-600" : "border-orange-200 bg-orange-50 text-orange-600"}`}>
                                                        {isAdmin(u) ? "Administrador" : "Bombero"}
                                                    </span>
                                                </td>
                                                <td className="py-2.5 px-3">
                                                    <span className={`inline-flex items-center px-3 py-1 rounded-full text-xs font-medium border ${u.is_active ? "border-green-200 bg-green-50 text-green-600" : "border-gray-200 bg-gray-50 text-gray-400"}`}>
                                                        {u.is_active ? "Activo" : "Inactivo"}
                                                    </span>
                                                </td>
                                                <td className="py-2.5 px-3 text-right">
                                                    <div className="flex gap-1 justify-end">
                                                        <Button variant="ghost" size="sm" onClick={() => openHistory(u)} title="Historial">
                                                            <Clock className="h-3.5 w-3.5" />
                                                        </Button>
                                                        <Button variant="ghost" size="sm" onClick={() => openEdit(u)}>
                                                            <Pencil className="h-3.5 w-3.5" />
                                                        </Button>
                                                        <Button variant="ghost" size="sm" onClick={() => toggleActive(u)}>
                                                            <Power className={`h-3.5 w-3.5 ${u.is_active ? "text-destructive" : "text-green-600"}`} />
                                                        </Button>
                                                    </div>
                                                </td>
                                            </tr>
                                        ))}
                                    </tbody>
                                </table>
                            </div>
                        )}
                    </CardContent>
                </Card>
            </div>

            {/* History Modal */}
            <Dialog open={historyOpen} onOpenChange={setHistoryOpen}>
                <DialogContent className="sm:max-w-2xl max-h-[80vh] overflow-hidden flex flex-col">
                    <DialogHeader>
                        <DialogTitle>Historial — {historyUser?.nombres} {historyUser?.apellidos}</DialogTitle>
                    </DialogHeader>
                    <div className="overflow-y-auto flex-1">
                        {historyLoading ? (
                            <p className="text-muted-foreground text-center py-8">Cargando...</p>
                        ) : history.length === 0 ? (
                            <p className="text-muted-foreground text-center py-8">Sin transacciones</p>
                        ) : (
                            <table className="w-full text-sm">
                                <thead className="sticky top-0 bg-background">
                                    <tr className="border-b">
                                        <th className="text-left py-2 px-3 font-medium text-muted-foreground">Fecha</th>
                                        <th className="text-left py-2 px-3 font-medium text-muted-foreground">Hora</th>
                                        <th className="text-left py-2 px-3 font-medium text-muted-foreground">Tipo</th>
                                        <th className="text-right py-2 px-3 font-medium text-muted-foreground">Galones</th>
                                        <th className="text-right py-2 px-3 font-medium text-muted-foreground">Puntos</th>
                                    </tr>
                                </thead>
                                <tbody>
                                    {history.map((tx) => {
                                        const d = new Date(tx.created_at)
                                        const isEarn = tx.transaction_type === "earn"
                                        return (
                                            <tr key={tx.id} className="border-b last:border-0 hover:bg-muted/50">
                                                <td className="py-2 px-3 whitespace-nowrap">
                                                    {d.toLocaleDateString("es-DO")}
                                                </td>
                                                <td className="py-2 px-3 whitespace-nowrap text-muted-foreground">
                                                    {d.toLocaleTimeString("es-DO", { hour: "2-digit", minute: "2-digit" })}
                                                </td>
                                                <td className="py-2 px-3">
                                                    <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium border ${isEarn ? "border-emerald-200 bg-emerald-50 text-emerald-600" : "border-red-200 bg-red-50 text-red-500"}`}>
                                                        {isEarn ? "Acumulación" : "Canje"}
                                                    </span>
                                                </td>
                                                <td className="py-2 px-3 text-right">
                                                    {tx.gallons_amount > 0 ? tx.gallons_amount.toFixed(1) : "-"}
                                                </td>
                                                <td className={`py-2 px-3 text-right font-medium ${isEarn ? "text-emerald-600" : "text-red-500"}`}>
                                                    {isEarn ? "+" : ""}{Math.abs(tx.points).toFixed(0)}
                                                </td>
                                            </tr>
                                        )
                                    })}
                                </tbody>
                            </table>
                        )}
                    </div>
                </DialogContent>
            </Dialog>
        </>
    )
}
