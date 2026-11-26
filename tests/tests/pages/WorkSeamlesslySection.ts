import {Locator, Page} from "@playwright/test";

export class WorkSeamlesslySection {
    readonly page: Page;
    readonly section: Locator;

    constructor(page: Page) {
        this.page = page;
        this.section = page.locator('section', {
            has: page.getByText('Works seamlessly on', {exact: true}),
        });
    }

    async goto() {
        await this.page.goto('/');
    }
}