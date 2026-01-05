import { test, expect } from '@playwright/test'

test.describe('Authentication', () => {
  test('shows login page when not authenticated', async ({ page }) => {
    await page.goto('/')
    await expect(page.getByRole('heading', { name: 'Login' })).toBeVisible()
    await expect(page.getByPlaceholder('Enter your email')).toBeVisible()
  })

  test('can login with email', async ({ page }) => {
    await page.goto('/')
    await page.getByPlaceholder('Enter your email').fill('test@example.com')
    await page.getByRole('button', { name: 'Login' }).click()

    // Should redirect to home page after login
    await expect(page.getByRole('heading', { name: 'Hungr' })).toBeVisible()
  })

  test('persists login across page reloads', async ({ page }) => {
    await page.goto('/')
    await page.getByPlaceholder('Enter your email').fill('test@example.com')
    await page.getByRole('button', { name: 'Login' }).click()

    await expect(page.getByRole('heading', { name: 'Hungr' })).toBeVisible()

    // Reload the page
    await page.reload()

    // Should still be logged in
    await expect(page.getByRole('heading', { name: 'Hungr' })).toBeVisible()
  })

  test('can logout', async ({ page }) => {
    // Login first
    await page.goto('/')
    await page.getByPlaceholder('Enter your email').fill('test@example.com')
    await page.getByRole('button', { name: 'Login' }).click()

    await expect(page.getByRole('heading', { name: 'Hungr' })).toBeVisible()

    // Click logout
    await page.getByRole('button', { name: 'Logout' }).click()

    // Should be back at login page
    await expect(page.getByRole('heading', { name: 'Login' })).toBeVisible()
  })
})
