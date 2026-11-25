import {Locator, Page} from "@playwright/test";

export class HowSection {
    readonly page: Page;
    readonly titleHowSection: Locator;
    readonly numberOneHowSection: Locator;
    readonly numberTwoHowSection: Locator;
    readonly numberThreeHowSection: Locator;


    constructor(page: Page) {
        this.page = page;
        this.titleHowSection = page.getByRole('heading', {name: 'Built for effortless browsing'});
        this.numberOneHowSection = page.locator('li', {
            has: page.getByText('Install the extension in Google Chrome, Firefox, or Edge.'),
        });

        this.numberTwoHowSection = page.locator('li', {
            has: page.getByText('Pick your preferred local currency.'),
        });

        this.numberThreeHowSection = page.locator('li', {
            has: page.getByText('Browse normally — prices convert automatically in place.'),
        });
    }

    async goto() {
        await this.page.goto('/');
    }
}