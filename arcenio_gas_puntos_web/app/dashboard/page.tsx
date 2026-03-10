"use client"

import { useEffect, useState } from "react"
import { apiGet } from "@/lib/api"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Fuel, TrendingUp, TrendingDown, Users, ArrowRightLeft } from "lucide-react"

interface TopClient {
  nombres: string
  apellidos: string
  cedula: string
  points: number
}

interface RecentTransaction {
  transaction_type: string
  client_name: string
  points: number
  gallons_amount: number
  created_at: string
}

interface DashboardData {
  total_gallons: number
  total_points_earned: number
  total_points_redeemed: number
  total_transactions: number
  total_clients: number
  month: number
  year: number
  top_clients: TopClient[] | null
  recent_transactions: RecentTransaction[] | null
}

const MONTH_NAMES = [
  "Enero", "Febrero", "Marzo", "Abril", "Mayo", "Junio",
  "Julio", "Agosto", "Septiembre", "Octubre", "Noviembre", "Diciembre",
]

export default function DashboardPage() {
  const [data, setData] = useState<DashboardData | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    async function fetchDashboard() {
      try {
        const d = await apiGet<DashboardData>("/reports/dashboard")
        setData(d)
      } catch {
        // silently fail
      } finally {
        setLoading(false)
      }
    }
    fetchDashboard()
  }, [])

  if (loading) {
    return (
      <div className="space-y-6">
        <div className="grid gap-4 md:grid-cols-5">
          {[1, 2, 3, 4, 5].map((i) => (
            <div key={i} className="bg-muted/50 h-24 rounded-xl animate-pulse" />
          ))}
        </div>
        <div className="grid gap-4 md:grid-cols-2">
          <div className="bg-muted/50 h-80 rounded-xl animate-pulse" />
          <div className="bg-muted/50 h-80 rounded-xl animate-pulse" />
        </div>
      </div>
    )
  }

  const kpis = [
    {
      label: "Galones",
      value: data?.total_gallons?.toFixed(1) ?? "0",
      suffix: "gal",
      icon: Fuel,
      color: "text-blue-600 bg-blue-50 border-blue-200",
    },
    {
      label: "Puntos Otorgados",
      value: data?.total_points_earned?.toFixed(0) ?? "0",
      suffix: "pts",
      icon: TrendingUp,
      color: "text-emerald-600 bg-emerald-50 border-emerald-200",
    },
    {
      label: "Puntos Canjeados",
      value: data?.total_points_redeemed?.toFixed(0) ?? "0",
      suffix: "pts",
      icon: TrendingDown,
      color: "text-red-500 bg-red-50 border-red-200",
    },
    {
      label: "Transacciones",
      value: data?.total_transactions?.toString() ?? "0",
      suffix: "",
      icon: ArrowRightLeft,
      color: "text-amber-600 bg-amber-50 border-amber-200",
    },
    {
      label: "Clientes",
      value: data?.total_clients?.toString() ?? "0",
      suffix: "",
      icon: Users,
      color: "text-purple-600 bg-purple-50 border-purple-200",
    },
  ]

  const topClients = data?.top_clients ?? []
  const recentTransactions = data?.recent_transactions ?? []
  const monthLabel = data ? `${MONTH_NAMES[(data.month - 1) % 12]} ${data.year}` : ""

  return (
    <div className="space-y-6">
      {/* Header */}
      <div>
        <h2 className="text-2xl font-bold tracking-tight">Dashboard</h2>
        <p className="text-sm text-muted-foreground mt-1">
          Resumen de actividad — {monthLabel}
        </p>
      </div>

      {/* KPI Row */}
      <div className="grid gap-4 grid-cols-2 md:grid-cols-5">
        {kpis.map((kpi) => (
          <Card key={kpi.label} className="relative overflow-hidden">
            <CardContent className="p-4">
              <div className="flex items-center justify-between mb-3">
                <span className="text-xs font-medium text-muted-foreground uppercase tracking-wider">
                  {kpi.label}
                </span>
                <div className={`p-1.5 rounded-lg border ${kpi.color}`}>
                  <kpi.icon className="h-3.5 w-3.5" />
                </div>
              </div>
              <div className="text-2xl font-bold tracking-tight">
                {kpi.value}
                {kpi.suffix && (
                  <span className="text-xs font-normal text-muted-foreground ml-1">
                    {kpi.suffix}
                  </span>
                )}
              </div>
            </CardContent>
          </Card>
        ))}
      </div>

      {/* Two column section */}
      <div className="grid gap-6 md:grid-cols-2">
        {/* Recent Transactions */}
        <Card>
          <CardHeader className="pb-3">
            <CardTitle className="text-base font-semibold">Últimas Transacciones</CardTitle>
            <p className="text-xs text-muted-foreground">Las 5 transacciones más recientes</p>
          </CardHeader>
          <CardContent>
            {recentTransactions.length === 0 ? (
              <p className="text-sm text-muted-foreground text-center py-6">Sin transacciones</p>
            ) : (
              <div className="space-y-3">
                {recentTransactions.map((tx, idx) => {
                  const isEarn = tx.transaction_type === "earn"
                  const d = new Date(tx.created_at)
                  return (
                    <div key={idx} className="flex items-center gap-3 py-1.5">
                      <div className={`flex-shrink-0 w-8 h-8 rounded-full flex items-center justify-center ${isEarn ? "bg-emerald-50 text-emerald-600" : "bg-red-50 text-red-500"}`}>
                        {isEarn ? <TrendingUp className="h-3.5 w-3.5" /> : <TrendingDown className="h-3.5 w-3.5" />}
                      </div>
                      <div className="flex-1 min-w-0">
                        <p className="text-sm font-medium truncate">{tx.client_name}</p>
                        <p className="text-xs text-muted-foreground">
                          {d.toLocaleDateString("es-DO")} · {d.toLocaleTimeString("es-DO", { hour: "2-digit", minute: "2-digit" })}
                          {isEarn && tx.gallons_amount > 0 ? ` · ${tx.gallons_amount.toFixed(1)} gal` : ""}
                        </p>
                      </div>
                      <div className={`text-sm font-semibold whitespace-nowrap ${isEarn ? "text-emerald-600" : "text-red-500"}`}>
                        {isEarn ? "+" : ""}{Math.abs(tx.points).toFixed(0)} pts
                      </div>
                    </div>
                  )
                })}
              </div>
            )}
          </CardContent>
        </Card>

        {/* Top Clients */}
        <Card>
          <CardHeader className="pb-3">
            <CardTitle className="text-base font-semibold">Top Clientes</CardTitle>
            <p className="text-xs text-muted-foreground">Mayor acumulación de puntos — {monthLabel}</p>
          </CardHeader>
          <CardContent>
            {topClients.length === 0 ? (
              <p className="text-sm text-muted-foreground text-center py-6">Sin datos</p>
            ) : (
              <div className="overflow-x-auto">
                <table className="w-full text-sm">
                  <thead>
                    <tr className="border-b">
                      <th className="text-left py-2 pr-2 font-medium text-muted-foreground text-xs">#</th>
                      <th className="text-left py-2 px-2 font-medium text-muted-foreground text-xs">CLIENTE</th>
                      <th className="text-right py-2 pl-2 font-medium text-muted-foreground text-xs">PUNTOS</th>
                    </tr>
                  </thead>
                  <tbody>
                    {topClients.map((client, idx) => (
                      <tr key={client.cedula} className="border-b last:border-0">
                        <td className="py-2.5 pr-2">
                          <span className={`inline-flex items-center justify-center w-6 h-6 rounded-full text-xs font-bold ${idx === 0 ? "bg-amber-100 text-amber-700" : idx === 1 ? "bg-gray-100 text-gray-600" : idx === 2 ? "bg-orange-100 text-orange-700" : "bg-muted text-muted-foreground"}`}>
                            {idx + 1}
                          </span>
                        </td>
                        <td className="py-2.5 px-2">
                          <p className="font-medium text-sm">{client.nombres} {client.apellidos}</p>
                          <p className="text-xs text-muted-foreground">{formatCedula(client.cedula)}</p>
                        </td>
                        <td className="py-2.5 pl-2 text-right">
                          <span className="font-semibold text-emerald-600">
                            {client.points.toFixed(0)}
                          </span>
                          <span className="text-xs text-muted-foreground ml-1">pts</span>
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
    </div>
  )
}

function formatCedula(raw: string): string {
  const d = raw.replace(/\D/g, "")
  if (d.length !== 11) return raw
  return `${d.substring(0, 3)}-${d.substring(3, 10)}-${d.substring(10)}`
}
