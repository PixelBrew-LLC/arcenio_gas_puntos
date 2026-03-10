"use client"

import { useEffect, useState, type FormEvent } from "react"
import { useParams } from "next/navigation"
import { apiGet, apiPost, apiPut } from "@/lib/api"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Button } from "@/components/ui/button"
import { Label } from "@/components/ui/label"
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogFooter,
    DialogHeader,
    DialogTitle,
} from "@/components/ui/dialog"
import { User, TrendingUp, TrendingDown, ArrowLeft, Pencil } from "lucide-react"
import Link from "next/link"
import { toast } from "sonner"

interface Client {
    id: string
    nombres: string
    apellidos: string
    cedula: string
    direccion?: string | null
    telefono?: string | null
}

interface Transaction {
    id: string
    transaction_type: string
    points: number
    gallons_amount: number
    processed_by_name?: string
    created_at: string
    expires_at?: string
}

function formatCedula(raw: string): string {
    const d = raw.replace(/\D/g, "")
    if (d.length !== 11) return raw
    return `${d.substring(0, 3)}-${d.substring(3, 10)}-${d.substring(10)}`
}

function formatTelefono(raw: string): string {
    if (!raw) return ""
    const d = raw.replace(/\D/g, "")
    if (d.length !== 10) return raw
    return `(${d.substring(0, 3)}) ${d.substring(3, 6)}-${d.substring(6)}`
}

function telefonoMask(value: string): string {
    const d = value.replace(/\D/g, "").substring(0, 10)
    let result = ""
    if (d.length > 0) result += "(" + d.substring(0, 3)
    if (d.length >= 4) result += ") " + d.substring(3, 6)
    if (d.length >= 7) result += "-" + d.substring(6, 10)
    return result
}

export default function ClientDetailPage() {
    const params = useParams()
    const cedula = params.cedula as string

    const [client, setClient] = useState<Client | null>(null)
    const [balance, setBalance] = useState(0)
    const [minRedeem, setMinRedeem] = useState(0)
    const [history, setHistory] = useState<Transaction[]>([])
    const [loading, setLoading] = useState(true)
    const [error, setError] = useState("")

    // Earn
    const [gallons, setGallons] = useState("")
    const [earnResult, setEarnResult] = useState<{ points_earned: number; new_balance: number } | null>(null)
    const [earnLoading, setEarnLoading] = useState(false)

    // Redeem
    const [redeemResult, setRedeemResult] = useState<{ points_redeemed: number; new_balance: number } | null>(null)
    const [redeemLoading, setRedeemLoading] = useState(false)
    const [confirmOpen, setConfirmOpen] = useState(false)

    // Edit
    const [editOpen, setEditOpen] = useState(false)
    const [editNombres, setEditNombres] = useState("")
    const [editApellidos, setEditApellidos] = useState("")
    const [editDireccion, setEditDireccion] = useState("")
    const [editTelefono, setEditTelefono] = useState("")
    const [editLoading, setEditLoading] = useState(false)
    const [editError, setEditError] = useState("")

    async function fetchClient() {
        try {
            const c = await apiGet<Client>(`/clients/${cedula}`)
            setClient(c)

            const b = await apiGet<{ balance: number; min_redeem: number }>(`/transactions/balance/${c.id}`)
            setBalance(b.balance)
            setMinRedeem(b.min_redeem)

            const h = await apiGet<Transaction[]>(`/transactions/history/${c.id}`)
            setHistory(h || [])
        } catch (err) {
            setError(err instanceof Error ? err.message : "Cliente no encontrado")
        } finally {
            setLoading(false)
        }
    }

    useEffect(() => { fetchClient() }, [cedula])

    async function handleEarn() {
        if (!client || !gallons) return
        setEarnLoading(true)
        try {
            const result = await apiPost<{ points_earned: number; new_balance: number }>("/transactions/earn", {
                client_id: client.id,
                gallons: parseFloat(gallons),
            })
            setEarnResult(result)
            setBalance(result.new_balance)
            setGallons("")
            const h = await apiGet<Transaction[]>(`/transactions/history/${client.id}`)
            setHistory(h || [])
            toast.success(`+${result.points_earned.toFixed(0)} puntos acumulados`)
        } catch (err) {
            const msg = err instanceof Error ? err.message : "Error al acumular"
            setError(msg)
            toast.error(msg)
        } finally {
            setEarnLoading(false)
        }
    }

    async function handleRedeem() {
        if (!client) return
        setRedeemLoading(true)
        try {
            const result = await apiPost<{ points_redeemed: number; new_balance: number }>("/transactions/redeem", {
                client_id: client.id,
            })
            setRedeemResult(result)
            setBalance(result.new_balance)
            const h = await apiGet<Transaction[]>(`/transactions/history/${client.id}`)
            setHistory(h || [])
            toast.success(`${result.points_redeemed.toFixed(0)} puntos canjeados`)
        } catch (err) {
            const msg = err instanceof Error ? err.message : "Error al canjear"
            setError(msg)
            toast.error(msg)
        } finally {
            setRedeemLoading(false)
        }
    }

    if (loading) return <p className="text-muted-foreground text-center py-8">Cargando...</p>

    function openEditModal() {
        if (!client) return
        setEditNombres(client.nombres)
        setEditApellidos(client.apellidos)
        setEditDireccion(client.direccion || "")
        setEditTelefono(client.telefono ? telefonoMask(client.telefono) : "")
        setEditError("")
        setEditOpen(true)
    }

    async function handleEdit(e: FormEvent) {
        e.preventDefault()
        if (!client) return
        setEditLoading(true)
        setEditError("")
        try {
            const body: Record<string, string> = {
                nombres: editNombres.trim(),
                apellidos: editApellidos.trim(),
            }
            if (editDireccion.trim()) body.direccion = editDireccion.trim()
            const rawTelefono = editTelefono.replace(/\D/g, "")
            if (rawTelefono) body.telefono = rawTelefono
            await apiPut(`/clients/${client.id}`, body)
            setEditOpen(false)
            setLoading(true)
            await fetchClient()
            toast.success("Cliente actualizado")
        } catch (err) {
            const msg = err instanceof Error ? err.message : "Error al actualizar"
            setEditError(msg)
            toast.error(msg)
        } finally {
            setEditLoading(false)
        }
    }
    if (error && !client) return (
        <div className="text-center py-8 space-y-3">
            <p className="text-destructive">{error}</p>
            <Button variant="outline" asChild><Link href="/dashboard/clients">Volver</Link></Button>
        </div>
    )

    return (
        <div className="space-y-6">
            <Button variant="ghost" size="sm" asChild>
                <Link href="/dashboard/clients">
                    <ArrowLeft className="mr-2 h-4 w-4" /> Volver a Clientes
                </Link>
            </Button>

            {client && (
                <>
                    {/* Client info */}
                    <Card>
                        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                            <div className="flex items-center gap-3">
                                <User className="h-5 w-5 text-primary" />
                                <CardTitle className="text-lg">{client.nombres} {client.apellidos}</CardTitle>
                            </div>
                            <Button variant="outline" size="sm" onClick={openEditModal}>
                                <Pencil className="mr-2 h-3.5 w-3.5" /> Editar
                            </Button>
                        </CardHeader>
                        <CardContent>
                            <div className="grid grid-cols-2 md:grid-cols-4 gap-3 text-sm">
                                <div>
                                    <p className="text-muted-foreground">Cédula</p>
                                    <p className="font-medium">{formatCedula(client.cedula)}</p>
                                </div>
                                <div>
                                    <p className="text-muted-foreground">Teléfono</p>
                                    <p className="font-medium">{client.telefono ? formatTelefono(client.telefono) : "-"}</p>
                                </div>
                                <div>
                                    <p className="text-muted-foreground">Dirección</p>
                                    <p className="font-medium">{client.direccion || "-"}</p>
                                </div>
                                <div>
                                    <p className="text-muted-foreground">Saldo</p>
                                    <p className="font-bold text-primary text-lg">{balance.toFixed(0)} pts</p>
                                </div>
                            </div>
                        </CardContent>
                    </Card>

                    {/* Operations */}
                    <div className="grid md:grid-cols-2 gap-4">
                        <Card>
                            <CardHeader className="pb-3">
                                <CardTitle className="text-base flex items-center gap-2">
                                    <TrendingUp className="h-4 w-4 text-green-600" /> Acumular Puntos
                                </CardTitle>
                            </CardHeader>
                            <CardContent className="space-y-3">
                                {earnResult ? (
                                    <div className="text-center space-y-2">
                                        <p className="text-green-600 font-bold text-lg">+{earnResult.points_earned.toFixed(0)} puntos</p>
                                        <p className="text-sm text-muted-foreground">Nuevo saldo: {earnResult.new_balance.toFixed(0)} pts</p>
                                        <Button variant="outline" size="sm" onClick={() => setEarnResult(null)}>Otra acumulación</Button>
                                    </div>
                                ) : (
                                    <>
                                        <div>
                                            <Label htmlFor="gallons">Galones</Label>
                                            <Input id="gallons" type="number" step="0.01" placeholder="Ej: 15.5" value={gallons} onChange={(e) => setGallons(e.target.value)} />
                                        </div>
                                        <Button onClick={handleEarn} disabled={earnLoading || !gallons} className="w-full">
                                            {earnLoading ? "Acumulando..." : "Acumular"}
                                        </Button>
                                    </>
                                )}
                            </CardContent>
                        </Card>

                        <Card>
                            <CardHeader className="pb-3">
                                <CardTitle className="text-base flex items-center gap-2">
                                    <TrendingDown className="h-4 w-4 text-primary" /> Canjear Puntos
                                </CardTitle>
                            </CardHeader>
                            <CardContent className="space-y-3">
                                {redeemResult ? (
                                    <div className="text-center space-y-2">
                                        <p className="text-primary font-bold text-lg">-{redeemResult.points_redeemed.toFixed(0)} puntos</p>
                                        <p className="text-sm text-muted-foreground">Nuevo saldo: {redeemResult.new_balance.toFixed(0)} pts</p>
                                        <Button variant="outline" size="sm" onClick={() => setRedeemResult(null)}>OK</Button>
                                    </div>
                                ) : (
                                    <>
                                        <p className="text-sm text-muted-foreground">Se canjean todos los puntos. Mínimo: {minRedeem.toFixed(0)} pts.</p>
                                        <p className="text-xl font-bold text-center">{balance.toFixed(0)} puntos</p>
                                        <Button onClick={() => setConfirmOpen(true)} disabled={redeemLoading || balance < minRedeem} variant="secondary" className="w-full">
                                            {redeemLoading ? "Canjeando..." : `Canjear ${balance.toFixed(0)} puntos`}
                                        </Button>

                                        <Dialog open={confirmOpen} onOpenChange={setConfirmOpen}>
                                            <DialogContent className="sm:max-w-sm">
                                                <DialogHeader>
                                                    <DialogTitle>Confirmar Canje</DialogTitle>
                                                    <DialogDescription>
                                                        ¿Estás seguro de que deseas canjear <strong>{balance.toFixed(0)} puntos</strong> del cliente <strong>{client?.nombres} {client?.apellidos}</strong>?
                                                    </DialogDescription>
                                                </DialogHeader>
                                                <DialogFooter className="gap-2 sm:gap-0">
                                                    <Button variant="outline" onClick={() => setConfirmOpen(false)}>Cancelar</Button>
                                                    <Button onClick={() => { setConfirmOpen(false); handleRedeem() }} disabled={redeemLoading}>
                                                        {redeemLoading ? "Canjeando..." : "Confirmar Canje"}
                                                    </Button>
                                                </DialogFooter>
                                            </DialogContent>
                                        </Dialog>
                                    </>
                                )}
                            </CardContent>
                        </Card>
                    </div>

                    {/* History */}
                    <Card>
                        <CardHeader>
                            <CardTitle className="text-base">Historial de Transacciones</CardTitle>
                        </CardHeader>
                        <CardContent>
                            {history.length === 0 ? (
                                <p className="text-sm text-muted-foreground text-center py-4">Sin transacciones</p>
                            ) : (
                                <div className="overflow-x-auto">
                                    <table className="w-full text-sm">
                                        <thead>
                                            <tr className="border-b">
                                                <th className="text-left py-2 px-3 font-medium text-muted-foreground">Fecha</th>
                                                <th className="text-left py-2 px-3 font-medium text-muted-foreground">Hora</th>
                                                <th className="text-left py-2 px-3 font-medium text-muted-foreground">Tipo</th>
                                                <th className="text-left py-2 px-3 font-medium text-muted-foreground">Usuario</th>
                                                <th className="text-left py-2 px-3 font-medium text-muted-foreground">Expiración</th>
                                                <th className="text-right py-2 px-3 font-medium text-muted-foreground">Galones</th>
                                                <th className="text-right py-2 px-3 font-medium text-muted-foreground">Puntos</th>
                                            </tr>
                                        </thead>
                                        <tbody>
                                            {history.map((tx) => {
                                                const isEarn = tx.transaction_type === "earn"
                                                const d = new Date(tx.created_at)
                                                return (
                                                    <tr key={tx.id} className="border-b last:border-0 hover:bg-muted/50">
                                                        <td className="py-2.5 px-3 whitespace-nowrap text-muted-foreground">
                                                            {d.toLocaleDateString("es-DO")}
                                                        </td>
                                                        <td className="py-2.5 px-3 whitespace-nowrap text-muted-foreground">
                                                            {d.toLocaleTimeString("es-DO", { hour: "2-digit", minute: "2-digit" })}
                                                        </td>
                                                        <td className="py-2.5 px-3">
                                                            <span className={`inline-flex items-center px-3 py-1 rounded-full text-xs font-medium border ${isEarn ? "border-blue-200 bg-blue-50 text-blue-600" : "border-red-200 bg-red-50 text-red-600"}`}>
                                                                {isEarn ? "Acumulación" : "Canje"}
                                                            </span>
                                                        </td>
                                                        <td className="py-2.5 px-3 text-muted-foreground">
                                                            {tx.processed_by_name || "-"}
                                                        </td>
                                                        <td className="py-2.5 px-3 whitespace-nowrap">
                                                            {tx.expires_at ? (() => {
                                                                const isExpired = new Date(tx.expires_at).getTime() < Date.now()
                                                                return (
                                                                    <span className={`inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium border ${!isExpired ? "border-green-200 bg-green-50 text-green-700" : "border-red-200 bg-red-50 text-red-700"}`}>
                                                                        {new Date(tx.expires_at).toLocaleDateString("es-DO")}
                                                                    </span>
                                                                )
                                                            })() : <span className="text-muted-foreground">-</span>}
                                                        </td>
                                                        <td className="py-2.5 px-3 text-right text-muted-foreground">
                                                            {isEarn && tx.gallons_amount > 0 ? tx.gallons_amount.toFixed(1) : "-"}
                                                        </td>
                                                        <td className={`py-2.5 px-3 text-right font-medium ${isEarn ? "text-blue-600" : "text-red-600"}`}>
                                                            {isEarn ? "+" : ""}{tx.points.toFixed(0)}
                                                        </td>
                                                    </tr>
                                                )
                                            })}
                                        </tbody>
                                    </table>
                                </div>
                            )}
                        </CardContent>
                    </Card>
                    {/* Edit Modal */}
                    <Dialog open={editOpen} onOpenChange={setEditOpen}>
                        <DialogContent className="sm:max-w-lg">
                            <DialogHeader>
                                <DialogTitle>Editar Cliente</DialogTitle>
                            </DialogHeader>
                            <form onSubmit={handleEdit} className="space-y-4 mt-2">
                                <div className="grid grid-cols-2 gap-4">
                                    <div>
                                        <Label>Nombres *</Label>
                                        <Input value={editNombres} onChange={(e) => setEditNombres(e.target.value)} required autoFocus />
                                    </div>
                                    <div>
                                        <Label>Apellidos *</Label>
                                        <Input value={editApellidos} onChange={(e) => setEditApellidos(e.target.value)} required />
                                    </div>
                                </div>
                                <div>
                                    <Label>Cédula</Label>
                                    <Input value={client ? formatCedula(client.cedula) : ""} disabled className="bg-muted" />
                                </div>
                                <div className="grid grid-cols-2 gap-4">
                                    <div>
                                        <Label>Dirección</Label>
                                        <Input value={editDireccion} onChange={(e) => setEditDireccion(e.target.value)} />
                                    </div>
                                    <div>
                                        <Label>Teléfono</Label>
                                        <Input placeholder="(###) ###-####" value={editTelefono} onChange={(e) => setEditTelefono(telefonoMask(e.target.value))} />
                                    </div>
                                </div>
                                {editError && <p className="text-destructive text-sm">{editError}</p>}
                                <Button type="submit" disabled={editLoading} className="w-full">
                                    {editLoading ? "Guardando..." : "Guardar Cambios"}
                                </Button>
                            </form>
                        </DialogContent>
                    </Dialog>
                </>
            )
            }
        </div >
    )
}
