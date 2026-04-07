import { ImageResponse } from 'next/og';

export const runtime = 'edge';

export async function GET() {
  return new ImageResponse(
    (
      <div
        style={{
          width: '1200px',
          height: '630px',
          display: 'flex',
          flexDirection: 'column',
          alignItems: 'center',
          justifyContent: 'center',
          background: 'linear-gradient(135deg, #0f172a 0%, #1e3a5f 50%, #0f2744 100%)',
          fontFamily: 'sans-serif',
        }}
      >
        <div
          style={{
            display: 'flex',
            flexDirection: 'column',
            alignItems: 'center',
            gap: '24px',
            padding: '48px',
          }}
        >
          <div
            style={{
              fontSize: '96px',
              fontWeight: '700',
              color: '#ffffff',
              letterSpacing: '-2px',
              lineHeight: 1,
            }}
          >
            FinPlan
          </div>
          <div
            style={{
              fontSize: '32px',
              fontWeight: '400',
              color: '#93c5fd',
              textAlign: 'center',
              maxWidth: '900px',
              lineHeight: 1.4,
            }}
          >
            将来の資産推移を可視化し、老後・緊急資金を計画するツール
          </div>
        </div>
      </div>
    ),
    {
      width: 1200,
      height: 630,
    }
  );
}
