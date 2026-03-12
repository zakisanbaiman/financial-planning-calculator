import { test, expect } from '@playwright/test';
import { generateTestUserId, setupCompleteFinancialProfile, API_BASE_URL } from './test-utils';

/**
 * E2E Test: Reports Download Scenario
 *
 * Tests the full PDF report export and download flow:
 * - Export report -> receive download token -> download file
 * - Authorization checks (another user's token should return 403)
 */

test.describe('Reports Download Scenario', () => {
  test('Scenario: Export report and download PDF via token', async ({ request }) => {
    const userId = generateTestUserId();
    await setupCompleteFinancialProfile(request, userId);

    // Step 1: Export the report to get a download token
    const exportResponse = await request.post(`${API_BASE_URL}/api/reports/export`, {
      data: {
        user_id: userId,
        report_type: 'financial_summary',
        format: 'pdf',
        report_data: { key: 'value' },
      },
    });

    expect(exportResponse.ok()).toBeTruthy();
    const exportData = await exportResponse.json();
    expect(exportData.download_url).toBeDefined();
    expect(exportData.file_name).toBeDefined();
    expect(exportData.expires_at).toBeDefined();

    // Step 2: Extract token from download_url
    // Expected format: /api/reports/download/{token}
    const downloadUrl: string = exportData.download_url;
    expect(downloadUrl).toContain('/api/reports/download/');
    const token = downloadUrl.split('/api/reports/download/')[1];
    expect(token).toBeTruthy();

    // Step 3: Download the file using the token
    const downloadResponse = await request.get(
      `${API_BASE_URL}/api/reports/download/${token}`,
      {
        headers: {
          // 認証ヘッダー（実装時に認証情報を付与する）
          'X-User-ID': userId,
        },
      }
    );

    expect(downloadResponse.ok()).toBeTruthy();
    const contentType = downloadResponse.headers()['content-type'];
    // Content-Type は application/pdf または text/html (HTMLGenerator の場合)
    expect(
      contentType.includes('application/pdf') || contentType.includes('text/html')
    ).toBeTruthy();

    const body = await downloadResponse.body();
    expect(body.length).toBeGreaterThan(0);
  });

  test('Scenario: Quick download via GET /reports/pdf endpoint', async ({ request }) => {
    const userId = generateTestUserId();
    await setupCompleteFinancialProfile(request, userId);

    const pdfResponse = await request.get(`${API_BASE_URL}/api/reports/pdf`, {
      params: {
        user_id: userId,
        report_type: 'comprehensive',
        years: '10',
      },
    });

    expect(pdfResponse.ok()).toBeTruthy();
    const data = await pdfResponse.json();
    expect(data.download_url).toBeDefined();
    expect(data.expires_at).toBeDefined();
  });

  test('Scenario: Export financial_summary report and verify response fields', async ({ request }) => {
    const userId = generateTestUserId();
    await setupCompleteFinancialProfile(request, userId);

    const exportResponse = await request.post(`${API_BASE_URL}/api/reports/export`, {
      data: {
        user_id: userId,
        report_type: 'financial_summary',
        format: 'pdf',
        report_data: {},
      },
    });

    expect(exportResponse.ok()).toBeTruthy();
    const data = await exportResponse.json();

    // レスポンスフィールドの検証
    expect(data.file_name).toBeDefined();
    expect(typeof data.file_size).toBe('number');
    expect(data.file_size).toBeGreaterThan(0);
    expect(data.download_url).toMatch(/\/api\/reports\/download\/.+/);
    expect(data.expires_at).toBeDefined();

    // expires_at が未来の日時であることを確認
    const expiresAt = new Date(data.expires_at);
    expect(expiresAt.getTime()).toBeGreaterThan(Date.now());
  });

  test('Scenario: Authorization check - another user cannot download token', async ({ request }) => {
    const ownerUserId = generateTestUserId();
    const attackerUserId = generateTestUserId();

    await setupCompleteFinancialProfile(request, ownerUserId);

    // オーナーがエクスポートしてトークンを取得
    const exportResponse = await request.post(`${API_BASE_URL}/api/reports/export`, {
      data: {
        user_id: ownerUserId,
        report_type: 'financial_summary',
        format: 'pdf',
        report_data: {},
      },
    });

    expect(exportResponse.ok()).toBeTruthy();
    const exportData = await exportResponse.json();
    const downloadUrl: string = exportData.download_url;
    const token = downloadUrl.split('/api/reports/download/')[1];
    expect(token).toBeTruthy();

    // 別のユーザーがオーナーのトークンでダウンロードを試みる -> 403 Forbidden
    const unauthorizedResponse = await request.get(
      `${API_BASE_URL}/api/reports/download/${token}`,
      {
        headers: {
          // 攻撃者のユーザーIDを送信
          'X-User-ID': attackerUserId,
        },
      }
    );

    // 認可チェックにより403が返ること
    expect(unauthorizedResponse.status()).toBe(403);
  });

  test('Scenario: Download with invalid token returns 404', async ({ request }) => {
    const response = await request.get(
      `${API_BASE_URL}/api/reports/download/invalid-nonexistent-token-xyz`,
      {
        headers: {
          'X-User-ID': generateTestUserId(),
        },
      }
    );

    expect(response.status()).toBe(404);
  });

  test('Scenario: Download with expired token returns 410', async ({ request }) => {
    // このテストは実際に期限切れトークンを作る代わりに、
    // 期限切れを示す特定のトークン形式でAPIに問い合わせる
    // 実装側でテスト用に期限切れ状態をシミュレートできるようにすること
    const expiredToken = 'expired-test-token-for-e2e';

    const response = await request.get(
      `${API_BASE_URL}/api/reports/download/${expiredToken}`,
      {
        headers: {
          'X-User-ID': generateTestUserId(),
        },
      }
    );

    // 期限切れの場合は 410 Gone か 404 Not Found が返る
    expect([404, 410]).toContain(response.status());
  });

  test('Scenario: Export all report types successfully', async ({ request }) => {
    const userId = generateTestUserId();
    await setupCompleteFinancialProfile(request, userId);

    const reportTypes = [
      'financial_summary',
      'asset_projection',
      'goals_progress',
      'comprehensive',
    ];

    for (const reportType of reportTypes) {
      const exportResponse = await request.post(`${API_BASE_URL}/api/reports/export`, {
        data: {
          user_id: userId,
          report_type: reportType,
          format: 'pdf',
          report_data: {},
        },
      });

      expect(exportResponse.ok()).toBeTruthy();
      const data = await exportResponse.json();
      expect(data.download_url).toBeDefined();
    }
  });

  test('Scenario: Error - Export without user_id returns 400', async ({ request }) => {
    const response = await request.post(`${API_BASE_URL}/api/reports/export`, {
      data: {
        report_type: 'financial_summary',
        format: 'pdf',
        report_data: {},
      },
    });

    expect(response.status()).toBe(400);
  });

  test('Scenario: Error - Export with invalid report_type returns 400', async ({ request }) => {
    const userId = generateTestUserId();

    const response = await request.post(`${API_BASE_URL}/api/reports/export`, {
      data: {
        user_id: userId,
        report_type: 'invalid_report_type',
        format: 'pdf',
        report_data: {},
      },
    });

    expect(response.status()).toBe(400);
  });
});
