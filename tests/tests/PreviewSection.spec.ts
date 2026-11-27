import {expect, test} from "@playwright/test";
import {PreviewSection} from "./pages/PreviewSection";

test.describe('Currency converter preview', () => {

    test.beforeEach(async ({page}) => {
        const previewSection = new PreviewSection(page);
        await previewSection.goto()

        await expect(previewSection.previewForm).toBeVisible();
    });

    test('Shows the preview label and default amounts', async ({page}) => {
        const previewSection = new PreviewSection(page);

        await expect(previewSection.previewForm.getByText('Preview', {exact: true})).toBeVisible();

        await expect(previewSection.previewForm.getByText('You send', {exact: true})).toBeVisible();
        await expect(previewSection.previewForm.getByText('100.00', {exact: true})).toBeVisible();

        await expect(previewSection.previewForm.getByText('They receive', {exact: true})).toBeVisible();
        await expect(previewSection.previewForm.getByText('—', {exact: true})).toBeVisible();
    });

    test('shows USD as the send currency and EUR as the receive currency', async ({page}) => {
        const previewSection = new PreviewSection(page);
        const pills = previewSection.previewForm.getByText(/^(USD|EUR)$/);

        await expect(pills.nth(0)).toHaveText('USD');
        await expect(pills.nth(1)).toHaveText('EUR');
    });

    test('swap button is accessible and responds to a click', async ({ page }) => {
        const previewSection = new PreviewSection(page);

        const swap = previewSection.previewForm.getByRole('button', { name: 'Swap currencies' });

        await expect(swap).toBeVisible();
        await expect(swap).toBeEnabled();

        const styleBefore = await swap.getAttribute('style');
        await swap.click();

        await expect(swap).not.toHaveAttribute('style', styleBefore ?? '');
    });

    test('ticker lists the supported currencies', async ({page}) => {
        const previewSection = new PreviewSection(page);
        const ticker = previewSection.previewForm.locator('.animate-ticker');
        const codes = await ticker.locator('span.tracking-widest').allTextContents();

        // @ts-ignore
        const unique = [...new Set(codes)];

        expect(unique).toEqual([
            'USD', 'EUR', 'GBP', 'JPY', 'CHF', 'CAD', 'AUD', 'CNY',
            'INR', 'BRL', 'MXN', 'KRW', 'SEK', 'PLN', 'TRY', 'ZAR',
        ]);
    });
});