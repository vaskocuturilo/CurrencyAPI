import {expect, test} from "@playwright/test";
import {HowSection} from "./pages/HowSection";


test.describe('How section', () => {

    test.beforeEach(async ({page}) => {
        const howSection = new HowSection(page);
        await howSection.goto()
    });

    test('Check Features Section', async ({page}) => {
        const howSection = new HowSection(page);

        await expect(howSection.titleHowSection).toBeVisible()
        await expect(howSection.numberOneHowSection.locator('p')).toHaveText('Install the extension in Google Chrome, Firefox, or Edge.')
        await expect(howSection.numberTwoHowSection.locator('p')).toHaveText('Pick your preferred local currency.')
        await expect(howSection.numberThreeHowSection.locator('p')).toHaveText('Browse normally — prices convert automatically in place.')
    });
});