import Image from "next/image"

import { LoginForm } from "@/components/login-form"

export default function LoginPage() {
  return (
    <div className="bg-muted flex min-h-svh flex-col items-center justify-center gap-6 p-6 md:p-10">
      <div className="flex w-full max-w-sm flex-col gap-6">
        <div className="flex justify-center">
          <Image src="/arcenio_logo.png" alt="Arcenio Gas" width={120} height={120} priority />
        </div>
        <LoginForm />
      </div>
    </div>
  )
}
