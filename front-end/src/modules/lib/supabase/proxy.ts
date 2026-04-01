import { createServerClient } from '@supabase/ssr'
import { NextResponse, type NextRequest } from 'next/server'
import { PUBLIC_ROUTES } from '../routes'

export async function updateSession(request: NextRequest) {
  let supabaseResponse = NextResponse.next({
    request,
  })

  const supabase = createServerClient(
    process.env.NEXT_PUBLIC_SUPABASE_URL!,
    process.env.NEXT_PUBLIC_SUPABASE_PUBLISHABLE_KEY!,
    {
      cookies: {
        getAll() {
          return request.cookies.getAll()
        },
        setAll(cookiesToSet, headers) {
          cookiesToSet.forEach(({ name, value }) => request.cookies.set(name, value))
          supabaseResponse = NextResponse.next({
            request,
          })
          cookiesToSet.forEach(({ name, value, options }) =>
            supabaseResponse.cookies.set(name, value, options)
          )
          Object.entries(headers).forEach(([key, value]) =>
            supabaseResponse.headers.set(key, value)
          )
        },
      },
    }
  )
  const { data: { user } } = await supabase.auth.getUser();

  console.log(user)
  const pathname = request.nextUrl.pathname;
  const isPublicRoute = PUBLIC_ROUTES.includes(pathname);

  // Verificar se usuário está logado e não confirmou o email
  // Verificar se usuário já esta na rota que esta tentando acessar
  if (user && !user.email_confirmed_at && pathname !== '/confirm-email') {
    const url = request.nextUrl.clone();
    url.pathname = '/confirm-email';
    return NextResponse.redirect(url);
  }

  // Verificar se o usuário está logado e está tentando acessar uma rota privada
  // Verificar se usuário já esta na rota que esta tentando acessar
  if (user && user.email_confirmed_at && pathname !== '/feed') {
    // Internamente verificar se ele tem as roles para acessar essa página
    const url = request.nextUrl.clone();
    url.pathname = '/feed';
    return NextResponse.redirect(url);
  }

  // Verificar se o usuário não está logado e está tentando acessar rota privada
  if (!user && !isPublicRoute) {
    const url = request.nextUrl.clone();
    url.pathname = '/login';
    return NextResponse.redirect(url);
  }


  return supabaseResponse
}