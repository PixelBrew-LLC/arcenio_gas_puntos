"use client"

import { useEffect, useState, type FormEvent } from "react"
import { apiGet, apiPut } from "@/lib/api"
import { Card, CardContent } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Button } from "@/components/ui/button"
import { Save, Fuel, Award, ArrowDownToLine, CalendarClock } from "lucide-react"
import { toast } from "sonner"

interface Setting {
    key: string
    value: string
}

const SETTING_META: Record<string, { label: string; description: string; icon: typeof Fuel; color: string }> = {
    points_per_gallon: {
        label: "Puntos por Galón",
        description: "Cantidad de puntos otorgados por cada galón despachado",
        icon: Fuel,
        color: "text-blue-600 bg-blue-50 border-blue-200",
    },
    min_gallons: {
        label: "Mínimo de Galones",
        description: "Cantidad mínima de galones para otorgar puntos",
        icon: ArrowDownToLine,
        color: "text-amber-600 bg-amber-50 border-amber-200",
    },
    min_redeem_points: {
        label: "Mínimo para Canjear",
        description: "Cantidad mínima de puntos acumulados para poder canjear",
        icon: Award,
        color: "text-emerald-600 bg-emerald-50 border-emerald-200",
    },
    points_expiry_months: {
        label: "Vigencia (meses)",
        description: "Meses de vigencia de los puntos antes de expirar",
        icon: CalendarClock,
        color: "text-purple-600 bg-purple-50 border-purple-200",
    },
}

export default function SettingsPage() {
    const [settings, setSettings] = useState<Setting[]>([])
    const [loading, setLoading] = useState(true)
    const [saving, setSaving] = useState(false)
    const [error, setError] = useState("")

    useEffect(() => {
        async function fetch() {
            try {
                const data = await apiGet<Setting[]>("/settings")
                setSettings(data || [])
            } catch {
                // silently fail
            } finally {
                setLoading(false)
            }
        }
        fetch()
    }, [])

    function updateValue(key: string, value: string) {
        setSettings((prev) =>
            prev.map((s) => (s.key === key ? { ...s, value } : s))
        )
    }

    async function handleSave(e: FormEvent) {
        e.preventDefault()
        setSaving(true)
        setError("")

        try {
            for (const s of settings) {
                await apiPut("/settings", { key: s.key, value: String(s.value) })
            }
            toast.success("Configuración guardada")
        } catch (err) {
            const msg = err instanceof Error ? err.message : "Error al guardar"
            setError(msg)
            toast.error(msg)
        } finally {
            setSaving(false)
        }
    }

    if (loading) {
        return (
            <div className="space-y-4 max-w-2xl">
                {[1, 2, 3, 4].map((i) => (
                    <div key={i} className="bg-muted/50 h-20 rounded-xl animate-pulse" />
                ))}
            </div>
        )
    }

    return (
        <div className="max-w-2xl space-y-6">
            <div>
                <h2 className="text-2xl font-bold tracking-tight">Configuración</h2>
                <p className="text-sm text-muted-foreground mt-1">
                    Variables del programa de lealtad
                </p>
            </div>

            <form onSubmit={handleSave} className="space-y-4">
                {settings.map((s) => {
                    const meta = SETTING_META[s.key]
                    if (!meta) return null
                    const Icon = meta.icon
                    return (
                        <Card key={s.key} className="transition-shadow hover:shadow-sm">
                            <CardContent className="p-4">
                                <div className="flex items-center gap-4">
                                    <div className={`flex-shrink-0 p-2.5 rounded-lg border ${meta.color}`}>
                                        <Icon className="h-5 w-5" />
                                    </div>
                                    <div className="flex-1 min-w-0">
                                        <p className="text-sm font-semibold">{meta.label}</p>
                                        <p className="text-xs text-muted-foreground mt-0.5">{meta.description}</p>
                                    </div>
                                    <div className="flex-shrink-0 w-28">
                                        <Input
                                            id={s.key}
                                            type="number"
                                            step="any"
                                            value={s.value}
                                            onChange={(e) => updateValue(s.key, e.target.value)}
                                            className="text-right font-semibold text-base h-10"
                                        />
                                    </div>
                                </div>
                            </CardContent>
                        </Card>
                    )
                })}

                {error && <p className="text-destructive text-sm">{error}</p>}

                <div className="pt-2">
                    <Button type="submit" disabled={saving} size="lg">
                        <Save className="mr-2 h-4 w-4" />
                        {saving ? "Guardando..." : "Guardar Configuración"}
                    </Button>
                </div>
            </form>
        </div>
    )
}
