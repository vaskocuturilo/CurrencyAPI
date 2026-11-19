import {test, expect} from '@playwright/test';

test.describe('Document head', () => {
    test.beforeEach(async ({page}) => {
        await page.goto('http://localhost:8080');
    });

    test('has title', async ({page}) => {
        await expect(page).toHaveTitle("Real-time Currency Converter Extension — Every price, in your currency");
    });

    test('get header links', async ({page}) => {
        const currencyConverter = page.getByRole('link', {name: 'Currency Converter'})

        await currencyConverter.click()

        await expect(page.getByRole('link', {name: 'Currency Converter'})).toBeVisible()

        await expect(page.getByRole('link', {name: 'Features'})).toBeVisible()
        await expect(page.getByRole('link', {name: 'How it works'})).toBeVisible()
        await expect(page.getByRole('link', {name: 'FAQ'})).toBeVisible()
        await expect(page.getByRole('link', {name: 'Add to browser'})).toBeVisible()
    });

});
