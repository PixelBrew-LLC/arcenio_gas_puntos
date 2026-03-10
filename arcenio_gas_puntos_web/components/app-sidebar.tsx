"use client"

import * as React from "react"
import {
  Fuel,
  LayoutDashboard,
  Settings2,
  Users,
  UserCog,
  FileText,
  Search,
} from "lucide-react"

import Image from "next/image"

import { NavMain } from "@/components/nav-main"
import { NavUser } from "@/components/nav-user"
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarRail,
} from "@/components/ui/sidebar"

const navItems = [
  {
    title: "Dashboard",
    url: "/dashboard",
    icon: LayoutDashboard,
    isActive: true,
  },
  {
    title: "Clientes",
    url: "/dashboard/clients",
    icon: Fuel,
  },
  {
    title: "Reportes",
    url: "#",
    icon: FileText,
    items: [
      { title: "Historial de Transacciones", url: "/dashboard/reports/transactions" },
    ],
  },
  {
    title: "Personal",
    url: "/dashboard/users",
    icon: Users,
  },
  {
    title: "Configuración",
    url: "/dashboard/settings",
    icon: Settings2,
  },
]

export function AppSidebar({ ...props }: React.ComponentProps<typeof Sidebar>) {
  const [user, setUser] = React.useState({ name: "", role: "" })

  React.useEffect(() => {
    try {
      const stored = localStorage.getItem("user")
      if (stored) {
        const parsed = JSON.parse(stored)
        setUser({
          name: `${parsed.nombres} ${parsed.apellidos}`,
          role: parsed.role || "",
        })
      }
    } catch { }
  }, [])

  return (
    <Sidebar collapsible="icon" {...props}>
      <SidebarHeader>
        <SidebarMenu>
          <SidebarMenuItem>
            <SidebarMenuButton size="lg" asChild>
              <a href="/dashboard" className="flex items-center gap-2">
                <div className="flex aspect-square size-8 items-center justify-center rounded-lg bg-white p-0.5 shadow-sm">
                  <Image src="/arcenio_logo.png" alt="Arcenio Gas Logo" width={28} height={28} className="object-contain" />
                </div>
                <div className="grid flex-1 text-left text-sm leading-tight">
                  <span className="truncate font-semibold">Arcenio Gas</span>
                  <span className="truncate text-xs">Fidelización</span>
                </div>
              </a>
            </SidebarMenuButton>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarHeader>
      <SidebarContent>
        <NavMain items={navItems} />
      </SidebarContent>
      <SidebarFooter>
        <NavUser user={{ name: user.name, email: user.role, avatar: "" }} />
      </SidebarFooter>
      <SidebarRail />
    </Sidebar>
  )
}
