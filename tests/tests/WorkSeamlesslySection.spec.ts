import {expect, test} from "@playwright/test";
import {WorkSeamlesslySection} from "./pages/WorkSeamlesslySection";

test.describe('Work seamlessly section', () => {

    test.beforeEach(async ({page}) => {
        const workSeamlesslySection = new WorkSeamlesslySection(page);
        await workSeamlesslySection.goto()
    });

    test('Checks title', async ({page}) => {
        const workSeamlesslySection = new WorkSeamlesslySection(page);

        await expect(workSeamlesslySection.section).toBeVisible();
        await expect(workSeamlesslySection.section.locator('p').first()).toHaveText('Works seamlessly on');


    });
    test('Checks badges', async ({page}) => {
        const workSeamlesslySection = new WorkSeamlesslySection(page);
        const badges = workSeamlesslySection.section.locator('div > span');
        await expect(badges).toHaveCount(3);
        await expect(badges).toHaveText(['Chrome', 'Firefox', 'Edge']);
        for (let i = 0; i < 3; i++) {
            await expect(badges.nth(i).locator('svg.lucide-globe')).toBeVisible();
        }
    });
});