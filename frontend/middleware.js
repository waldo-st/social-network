import { NextResponse } from 'next/server';

export function middleware(req) {
  let token = req.cookies.get('token');
  
  if (!token && req.nextUrl.pathname !== '/login') {
    return NextResponse.redirect(new URL('/login', req.url))
  }
  
}

export const config = {
  matcher: '/',
}