import {Locator, Page} from "@playwright/test";

export class PreviewSection {
    readonly page: Page;
    readonly previewForm: Locator;


    constructor(page: Page) {
        this.page = page;
        this.previewForm = page
            .locator('div')
            .filter({ has: page.getByText('You send', { exact: true }) })
            .filter({ has: page.getByRole('button', { name: 'Swap currencies' }) })
            .last();
    }

    async goto() {
        await this.page.goto('/');
    }
}