import { test, expect } from '@playwright/test';

const BASE_URL = 'http://localhost:3000';

test('ログインフローの確認', async ({ page }) => {
  // ログインページへアクセス
  await page.goto(`${BASE_URL}/login`);
  await expect(page).toHaveURL(/login/);
  
  // テストユーザーを事前にAPIで作成
  const uniqueId = Date.now();
  const email = `e2etest-${uniqueId}@example.com`;
  const password = 'TestPass123!';
  
  const res = await page.request.post('http://localhost:8080/api/auth/register', {
    data: { email, password }
  });
  console.log('Register status:', res.status());

  // オンボーディングモーダルが表示されていたらスキップ
  const skipButton = page.getByText('スキップ');
  if (await skipButton.isVisible({ timeout: 3000 }).catch(() => false)) {
    await skipButton.click();
    await page.goto(`${BASE_URL}/login`);
  }

  // isLoading が解除されるまで待つ
  await page.waitForFunction(() => {
    const btn = document.querySelector('button[type="submit"]') as HTMLButtonElement;
    return btn && !btn.disabled;
  }, { timeout: 15000 });

  // フォームに入力
  await page.fill('input[type="email"]', email);
  await page.fill('input[type="password"]', password);
  
  // ログインボタンをクリック
  await page.click('button[type="submit"]');
  
  // ダッシュボードへ遷移することを確認
  await page.waitForURL(/dashboard/, { timeout: 10000 });
  console.log('遷移先URL:', page.url());
  
  await expect(page).toHaveURL(/dashboard/);
});
