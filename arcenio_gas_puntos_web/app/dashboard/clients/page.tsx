"use client"

import { useEffect, useState, type FormEvent } from "react"
import { useRouter } from "next/navigation"
import { apiGet, apiPost } from "@/lib/api"
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
import { UserPlus, ChevronLeft, ChevronRight, Search } from "lucide-react"

interface Client {
    id: string
    nombres: string
    apellidos: string
    cedula: string
    direccion?: string | null
    telefono?: string | null
    points_balance: number
    created_at: string
}

function formatCedula(raw: string): string {
    const d = raw.replace(/\D/g, "")
    if (d.length !== 11) return raw
    return `${d.substring(0, 3)}-${d.substring(3, 10)}-${d.substring(10)}`
}

function cedulaMask(value: string): string {
    const d = value.replace(/\D/g, "").substring(0, 11)
    let result = ""
    for (let i = 0; i < d.length; i++) {
        if (i === 3 || i === 10) result += "-"
        result += d[i]
    }
    return result
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

const PAGE_SIZE = 10

export default function ClientsPage() {
    const router = useRouter()
    const [clients, setClients] = useState<Client[]>([])
    const [loading, setLoading] = useState(true)
    const [search, setSearch] = useState("")
    const [filterType, setFilterType] = useState("all") // 'all' | 'cedula' | 'nombre' | 'apellido'
    const [page, setPage] = useState(1)
    const [dialogOpen, setDialogOpen] = useState(false)

    // Form
    const [nombres, setNombres] = useState("")
    const [apellidos, setApellidos] = useState("")
    const [cedula, setCedula] = useState("")
    const [direccion, setDireccion] = useState("")
    const [telefono, setTelefono] = useState("")
    const [formLoading, setFormLoading] = useState(false)
    const [formError, setFormError] = useState("")

    async function fetchClients() {
        try {
            const data = await apiGet<Client[]>("/clients")
            setClients(data || [])
        } catch {
            // silently fail
        } finally {
            setLoading(false)
        }
    }

    useEffect(() => { fetchClients() }, [])

    // Filter
    const filtered = clients.filter((c) => {
        if (!search) return true
        const q = search.toLowerCase()
        const rawSearch = search.replace(/\D/g, "")

        if (filterType === "cedula") return c.cedula.includes(rawSearch)
        if (filterType === "nombre") return c.nombres.toLowerCase().includes(q)
        if (filterType === "apellido") return c.apellidos.toLowerCase().includes(q)

        // 'all'
        return (
            c.nombres.toLowerCase().includes(q) ||
            c.apellidos.toLowerCase().includes(q) ||
            c.cedula.includes(rawSearch) ||
            (c.telefono && c.telefono.includes(rawSearch))
        )
    })

    // Pagination
    const totalPages = Math.max(1, Math.ceil(filtered.length / PAGE_SIZE))
    const paginated = filtered.slice((page - 1) * PAGE_SIZE, page * PAGE_SIZE)

    // Reset page when search changes
    useEffect(() => { setPage(1) }, [search])

    function handleSearchChange(val: string) {
        if (filterType === "cedula") {
            setSearch(cedulaMask(val))
            return
        }

        if (filterType === "all") {
            const isMostlyNumbers = /^[0-9-]+$/.test(val)
            if (isMostlyNumbers && val.replace(/\D/g, "").length <= 11) {
                setSearch(cedulaMask(val))
                return
            }
        }

        setSearch(val)
    }

    // Effect to reformat content if filter changes and it was looking at cedula
    useEffect(() => {
        setSearch("")
    }, [filterType])

    function resetForm() {
        setNombres(""); setApellidos(""); setCedula(""); setDireccion(""); setTelefono("")
        setFormError("")
    }

    async function handleCreate(e: FormEvent) {
        e.preventDefault()
        const rawCedula = cedula.replace(/\D/g, "")
        if (rawCedula.length !== 11) {
            setFormError("La cédula debe tener 11 dígitos")
            return
        }

        setFormLoading(true)
        setFormError("")

        try {
            const body: Record<string, string> = {
                nombres: nombres.trim(),
                apellidos: apellidos.trim(),
                cedula: rawCedula,
            }
            if (direccion.trim()) body.direccion = direccion.trim()
            const rawTelefono = telefono.replace(/\D/g, "")
            if (rawTelefono) body.telefono = rawTelefono

            await apiPost("/clients", body)
            setDialogOpen(false)
            resetForm()
            setLoading(true)
            await fetchClients()
        } catch (err) {
            setFormError(err instanceof Error ? err.message : "Error al registrar")
        } finally {
            setFormLoading(false)
        }
    }

    return (
        <div className="space-y-6">
            <div className="flex items-center justify-between">
                <h2 className="text-2xl font-bold">Clientes</h2>
                <Dialog open={dialogOpen} onOpenChange={(open) => { setDialogOpen(open); if (!open) resetForm() }}>
                    <DialogTrigger asChild>
                        <Button>
                            <UserPlus className="mr-2 h-4 w-4" />
                            Nuevo Cliente
                        </Button>
                    </DialogTrigger>
                    <DialogContent className="sm:max-w-lg">
                        <DialogHeader>
                            <DialogTitle>Registrar Cliente</DialogTitle>
                        </DialogHeader>
                        <form onSubmit={handleCreate} className="space-y-4 mt-2">
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
                            <div>
                                <Label>Cédula *</Label>
                                <Input placeholder="###-#######-#" value={cedula} onChange={(e) => setCedula(cedulaMask(e.target.value))} required />
                            </div>
                            <div className="grid grid-cols-2 gap-4">
                                <div>
                                    <Label>Dirección</Label>
                                    <Input value={direccion} onChange={(e) => setDireccion(e.target.value)} />
                                </div>
                                <div>
                                    <Label>Teléfono</Label>
                                    <Input placeholder="(###) ###-####" value={telefono} onChange={(e) => setTelefono(telefonoMask(e.target.value))} />
                                </div>
                            </div>
                            {formError && <p className="text-destructive text-sm">{formError}</p>}
                            <Button type="submit" disabled={formLoading} className="w-full">
                                {formLoading ? "Registrando..." : "Registrar Cliente"}
                            </Button>
                        </form>
                    </DialogContent>
                </Dialog>
            </div>

            {/* Search */}
            <div className="flex flex-col sm:flex-row gap-3">
                <div className="w-full sm:w-[200px]">
                    <Select value={filterType} onValueChange={setFilterType}>
                        <SelectTrigger>
                            <SelectValue placeholder="Filtrar por..." />
                        </SelectTrigger>
                        <SelectContent>
                            <SelectItem value="all">Todos los campos</SelectItem>
                            <SelectItem value="cedula">Cédula</SelectItem>
                            <SelectItem value="nombre">Nombres</SelectItem>
                            <SelectItem value="apellido">Apellidos</SelectItem>
                        </SelectContent>
                    </Select>
                </div>
                <div className="relative flex-1 max-w-sm">
                    <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
                    <Input
                        placeholder={
                            filterType === "cedula" ? "Escribe la cédula..." :
                                filterType === "nombre" ? "Buscar por nombres..." :
                                    filterType === "apellido" ? "Buscar por apellidos..." :
                                        "Buscar por cédula, nombre o apellido..."
                        }
                        value={search}
                        onChange={(e) => handleSearchChange(e.target.value)}
                        className="pl-9"
                    />
                </div>
            </div>

            {/* Table */}
            <Card>
                <CardContent className="pt-6">
                    {loading ? (
                        <p className="text-muted-foreground text-center py-8">Cargando...</p>
                    ) : filtered.length === 0 ? (
                        <p className="text-muted-foreground text-center py-8">
                            {search ? "No se encontraron resultados" : "No hay clientes registrados"}
                        </p>
                    ) : (
                        <>
                            <div className="overflow-x-auto">
                                <table className="w-full text-sm">
                                    <thead>
                                        <tr className="border-b">
                                            <th className="text-left py-2 px-3 font-medium text-muted-foreground w-32">Cédula</th>
                                            <th className="text-left py-2 px-3 font-medium text-muted-foreground">Nombres</th>
                                            <th className="text-left py-2 px-3 font-medium text-muted-foreground">Apellidos</th>
                                            <th className="text-left py-2 px-3 font-medium text-muted-foreground hidden md:table-cell">Teléfono</th>
                                            <th className="text-left py-2 px-3 font-medium text-muted-foreground hidden lg:table-cell">Dirección</th>
                                            <th className="text-left py-2 px-3 font-medium text-muted-foreground hidden md:table-cell">Registro</th>
                                            <th className="text-right py-2 px-3 font-medium text-muted-foreground">Puntos</th>
                                        </tr>
                                    </thead>
                                    <tbody>
                                        {paginated.map((c) => (
                                            <tr
                                                key={c.id}
                                                className="border-b last:border-0 hover:bg-muted/50 cursor-pointer"
                                                onClick={() => router.push(`/dashboard/clients/${c.cedula}`)}
                                            >
                                                <td className="py-2.5 px-3 whitespace-nowrap">{formatCedula(c.cedula)}</td>
                                                <td className="py-2.5 px-3 font-medium">{c.nombres}</td>
                                                <td className="py-2.5 px-3">{c.apellidos}</td>
                                                <td className="py-2.5 px-3 hidden md:table-cell text-muted-foreground whitespace-nowrap">{c.telefono ? formatTelefono(c.telefono) : "-"}</td>
                                                <td className="py-2.5 px-3 hidden lg:table-cell text-muted-foreground">{c.direccion || "-"}</td>
                                                <td className="py-2.5 px-3 hidden md:table-cell text-muted-foreground whitespace-nowrap">
                                                    {new Date(c.created_at).toLocaleDateString("es-DO")}
                                                </td>
                                                <td className="py-2.5 px-3 text-right font-medium text-primary whitespace-nowrap">
                                                    {c.points_balance ? c.points_balance.toFixed(0) : "0"}
                                                </td>
                                            </tr>
                                        ))}
                                    </tbody>
                                </table>
                            </div>

                            {/* Pagination */}
                            <div className="flex items-center justify-between pt-4 border-t mt-4">
                                <p className="text-sm text-muted-foreground">
                                    {filtered.length} cliente{filtered.length !== 1 ? "s" : ""}
                                </p>
                                <div className="flex items-center gap-2">
                                    <Button
                                        variant="outline"
                                        size="sm"
                                        disabled={page <= 1}
                                        onClick={() => setPage(page - 1)}
                                    >
                                        <ChevronLeft className="h-4 w-4" />
                                    </Button>
                                    <span className="text-sm text-muted-foreground">
                                        {page} / {totalPages}
                                    </span>
                                    <Button
                                        variant="outline"
                                        size="sm"
                                        disabled={page >= totalPages}
                                        onClick={() => setPage(page + 1)}
                                    >
                                        <ChevronRight className="h-4 w-4" />
                                    </Button>
                                </div>
                            </div>
                        </>
                    )}
                </CardContent>
            </Card>
        </div>
    )
}
