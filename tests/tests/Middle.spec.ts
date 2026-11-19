import {test, expect} from '@playwright/test';

test.describe('Document head', () => {
    test.beforeEach(async ({page}) => {
        await page.goto('/');
    });

    test('get middle', async ({page}) => {
        const currencyConverter = page.getByRole('link', {name: 'Currency Converter'})

        await currencyConverter.click()

        await expect(page.getByRole('link', {name: 'Currency Converter'})).toBeVisible()

        await expect(page.getByRole('link', {name: 'Features'})).toBeVisible()
        await expect(page.getByRole('link', {name: 'How it works'})).toBeVisible()
        await expect(page.getByRole('link', {name: 'FAQ'})).toBeVisible()
        await expect(page.getByRole('link', {name: 'Add to browser'})).toBeVisible()
    });

});
