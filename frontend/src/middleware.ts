import { NextRequest, NextResponse } from 'next/server';

export function middleware(request: NextRequest) {
  const backendUrl = process.env.BACKEND_URL;

  if (!backendUrl) {
    // BACKEND_URL未設定の場合（ローカル開発時）はリライトしない
    // next.config.js の rewrites が localhost:8080 にフォールバック
    return NextResponse.next();
  }

  // /api/* リクエストをバックエンドにリライト
  const url = new URL(request.nextUrl.pathname + request.nextUrl.search, backendUrl);
  return NextResponse.rewrite(url);
}

export const config = {
  matcher: '/api/:path*',
};
