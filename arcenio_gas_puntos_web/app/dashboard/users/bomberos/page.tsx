"use client"

import { useEffect, useState, type FormEvent } from "react"
import { apiGet, apiPost, apiPut, apiPatch } from "@/lib/api"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Button } from "@/components/ui/button"
import { Label } from "@/components/ui/label"
import { UserPlus, Pencil, Power } from "lucide-react"

interface User {
    id: string
    nombres: string
    apellidos: string
    cedula: string
    telefono: string
    direccion: string
    username: string
    is_active: boolean
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

export default function BomberosPage() {
    const [users, setUsers] = useState<User[]>([])
    const [loading, setLoading] = useState(true)
    const [showForm, setShowForm] = useState(false)
    const [editing, setEditing] = useState<User | null>(null)
    const [error, setError] = useState("")

    // Form fields
    const [nombres, setNombres] = useState("")
    const [apellidos, setApellidos] = useState("")
    const [cedula, setCedula] = useState("")
    const [telefono, setTelefono] = useState("")
    const [direccion, setDireccion] = useState("")
    const [username, setUsername] = useState("")
    const [clave, setClave] = useState("")
    const [formLoading, setFormLoading] = useState(false)

    async function fetchUsers() {
        try {
            const data = await apiGet<User[]>("/users/bomberos")
            setUsers(data || [])
        } catch {
            // silently fail
        } finally {
            setLoading(false)
        }
    }

    useEffect(() => { fetchUsers() }, [])

    function resetForm() {
        setNombres("")
        setApellidos("")
        setCedula("")
        setTelefono("")
        setDireccion("")
        setUsername("")
        setClave("")
        setEditing(null)
        setError("")
    }

    function openEdit(user: User) {
        setEditing(user)
        setNombres(user.nombres)
        setApellidos(user.apellidos)
        setCedula(formatCedula(user.cedula))
        setTelefono(user.telefono)
        setDireccion(user.direccion)
        setUsername(user.username)
        setClave("")
        setShowForm(true)
    }

    async function handleSubmit(e: FormEvent) {
        e.preventDefault()
        setFormLoading(true)
        setError("")

        const body: Record<string, string> = {
            nombres: nombres.trim(),
            apellidos: apellidos.trim(),
            cedula: cedula.replace(/\D/g, ""),
            telefono: telefono.trim(),
            direccion: direccion.trim(),
            username: username.trim(),
            role: "bombero",
        }
        if (clave) body.clave = clave

        try {
            if (editing) {
                await apiPut(`/users/bomberos/${editing.id}`, body)
            } else {
                if (!clave) { setError("La clave es requerida"); setFormLoading(false); return }
                body.clave = clave
                await apiPost("/users/bomberos", body)
            }
            setShowForm(false)
            resetForm()
            await fetchUsers()
        } catch (err) {
            setError(err instanceof Error ? err.message : "Error")
        } finally {
            setFormLoading(false)
        }
    }

    async function toggleActive(user: User) {
        try {
            await apiPatch(`/users/bomberos/${user.id}/active`, { is_active: !user.is_active })
            await fetchUsers()
        } catch {
            // silently fail
        }
    }

    return (
        <div className="space-y-6">
            <div className="flex items-center justify-between">
                <h2 className="text-2xl font-bold">Bomberos</h2>
                <Button onClick={() => { resetForm(); setShowForm(!showForm) }}>
                    <UserPlus className="mr-2 h-4 w-4" />
                    {showForm ? "Cancelar" : "Nuevo Bombero"}
                </Button>
            </div>

            {/* Form */}
            {showForm && (
                <Card>
                    <CardHeader>
                        <CardTitle className="text-base">{editing ? "Editar Bombero" : "Nuevo Bombero"}</CardTitle>
                    </CardHeader>
                    <CardContent>
                        <form onSubmit={handleSubmit} className="grid md:grid-cols-2 gap-4">
                            <div>
                                <Label>Nombres *</Label>
                                <Input value={nombres} onChange={(e) => setNombres(e.target.value)} required />
                            </div>
                            <div>
                                <Label>Apellidos *</Label>
                                <Input value={apellidos} onChange={(e) => setApellidos(e.target.value)} required />
                            </div>
                            <div>
                                <Label>Cédula *</Label>
                                <Input placeholder="###-#######-#" value={cedula} onChange={(e) => setCedula(formatCedula(e.target.value))} required />
                            </div>
                            <div>
                                <Label>Teléfono *</Label>
                                <Input value={telefono} onChange={(e) => setTelefono(e.target.value)} required />
                            </div>
                            <div>
                                <Label>Dirección *</Label>
                                <Input value={direccion} onChange={(e) => setDireccion(e.target.value)} required />
                            </div>
                            <div>
                                <Label>Usuario *</Label>
                                <Input value={username} onChange={(e) => setUsername(e.target.value)} required />
                            </div>
                            <div>
                                <Label>Clave (PIN numérico) {editing ? "(dejar vacío para no cambiar)" : "*"}</Label>
                                <Input
                                    type="password"
                                    inputMode="numeric"
                                    pattern="[0-9]*"
                                    value={clave}
                                    onChange={(e) => setClave(e.target.value.replace(/\D/g, ""))}
                                    required={!editing}
                                />
                            </div>
                            <div className="flex items-end">
                                <Button type="submit" disabled={formLoading} className="w-full">
                                    {formLoading ? "Guardando..." : editing ? "Actualizar" : "Crear Bombero"}
                                </Button>
                            </div>
                            {error && <p className="text-destructive text-sm col-span-2">{error}</p>}
                        </form>
                    </CardContent>
                </Card>
            )}

            {/* List */}
            <Card>
                <CardContent className="pt-6">
                    {loading ? (
                        <p className="text-muted-foreground text-center py-8">Cargando...</p>
                    ) : users.length === 0 ? (
                        <p className="text-muted-foreground text-center py-8">No hay bomberos registrados</p>
                    ) : (
                        <div className="overflow-x-auto">
                            <table className="w-full text-sm">
                                <thead>
                                    <tr className="border-b">
                                        <th className="text-left py-2 px-2 font-medium text-muted-foreground">Nombre</th>
                                        <th className="text-left py-2 px-2 font-medium text-muted-foreground">Cédula</th>
                                        <th className="text-left py-2 px-2 font-medium text-muted-foreground">Usuario</th>
                                        <th className="text-left py-2 px-2 font-medium text-muted-foreground">Estado</th>
                                        <th className="text-right py-2 px-2 font-medium text-muted-foreground">Acciones</th>
                                    </tr>
                                </thead>
                                <tbody>
                                    {users.map((u) => (
                                        <tr key={u.id} className="border-b last:border-0 hover:bg-muted/50">
                                            <td className="py-2 px-2">{u.nombres} {u.apellidos}</td>
                                            <td className="py-2 px-2">{formatCedula(u.cedula)}</td>
                                            <td className="py-2 px-2">{u.username}</td>
                                            <td className="py-2 px-2">
                                                <span className={`inline-flex items-center px-3 py-1 rounded-full text-xs font-medium border ${u.is_active ? "border-green-200 bg-green-50 text-green-600" : "border-gray-200 bg-gray-50 text-gray-400"}`}>
                                                    {u.is_active ? "Activo" : "Inactivo"}
                                                </span>
                                            </td>
                                            <td className="py-2 px-2 text-right">
                                                <div className="flex gap-1 justify-end">
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
    )
}
