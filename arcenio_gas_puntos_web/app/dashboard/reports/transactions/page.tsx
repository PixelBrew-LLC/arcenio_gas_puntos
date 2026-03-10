"use client"

import { useEffect, useState } from "react"
import { apiGet } from "@/lib/api"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"

interface Transaction {
    id: string
    client_id: string
    client_name?: string
    client_cedula?: string
    transaction_type: string
    points: number
    gallons_amount: number
    processed_by_name?: string
    created_at: string
    expires_at?: string | null
}

export default function TransactionReportsPage() {
    const [transactions, setTransactions] = useState<Transaction[]>([])
    const [loading, setLoading] = useState(true)

    useEffect(() => {
        async function fetch() {
            try {
                const data = await apiGet<Transaction[]>("/reports/transactions")
                setTransactions(data || [])
            } catch {
                // silently fail
            } finally {
                setLoading(false)
            }
        }
        fetch()
    }, [])

    return (
        <div className="space-y-6">
            <h2 className="text-2xl font-bold">Historial de Transacciones</h2>

            <Card>
                <CardHeader>
                    <CardTitle className="text-base">
                        {loading ? "Cargando..." : `${transactions.length} transacciones`}
                    </CardTitle>
                </CardHeader>
                <CardContent>
                    {!loading && transactions.length === 0 ? (
                        <p className="text-sm text-muted-foreground text-center py-4">Sin transacciones</p>
                    ) : !loading ? (
                        <div className="overflow-x-auto">
                            <table className="w-full text-sm">
                                <thead>
                                    <tr className="border-b">
                                        <th className="text-left py-2 px-3 font-medium text-muted-foreground">Fecha</th>
                                        <th className="text-left py-2 px-3 font-medium text-muted-foreground">Hora</th>
                                        <th className="text-left py-2 px-3 font-medium text-muted-foreground">Tipo</th>
                                        <th className="text-left py-2 px-3 font-medium text-muted-foreground">Cliente</th>
                                        <th className="text-left py-2 px-3 font-medium text-muted-foreground">Usuario</th>
                                        <th className="text-left py-2 px-3 font-medium text-muted-foreground">Expiración</th>
                                        <th className="text-right py-2 px-3 font-medium text-muted-foreground">Galones</th>
                                        <th className="text-right py-2 px-3 font-medium text-muted-foreground">Puntos</th>
                                    </tr>
                                </thead>
                                <tbody>
                                    {transactions.map((tx) => {
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
                                                <td className="py-2.5 px-3">
                                                    {tx.client_name || tx.client_cedula || tx.client_id.substring(0, 8)}
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
                    ) : null}
                </CardContent>
            </Card>
        </div>
    )
}
