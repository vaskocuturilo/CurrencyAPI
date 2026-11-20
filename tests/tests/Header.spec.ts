import {test, expect} from '@playwright/test';
import {HeaderPage} from "./pages/HeaderPage";

test.describe('Document head', () => {

    test.beforeEach(async ({page}) => {
        const headerPage = new HeaderPage(page);
        await headerPage.goto()
    });

    test('get middle', async ({page}) => {
        const headerPage = new HeaderPage(page);

        await headerPage.checkCurrencyConverterButton()

        await expect(headerPage.getFeaturesLink).toBeVisible()
        await expect(headerPage.getHowItWorksLink).toBeVisible()
        await expect(headerPage.getFAQLink).toBeVisible()
        await expect(headerPage.getAddToBrowserLink).toBeVisible()
    });
});
