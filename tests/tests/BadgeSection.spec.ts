import {expect, test} from "@playwright/test";
import {BadgesSection} from "./pages/BadgesSection";

test.describe('Badges section', () => {

    test.beforeEach(async ({page}) => {
        const badgesSection = new BadgesSection(page);
        await badgesSection.goto()
    });

    test('Checks badges', async ({page}) => {
        const badgesSection = new BadgesSection(page);

        await expect(badgesSection.getBadge).toHaveText("Add to Chrome")
        await expect(badgesSection.getBadge).toHaveText("Add to Firefox")
        await expect(badgesSection.getBadge).toHaveText("Add to Edge")

    });
});