import { act } from 'react'
import { createRoot } from 'react-dom/client'
import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it, vi } from 'vitest'

import { Tabs, TabsContent, TabsList, TabsTrigger } from './tabs'

describe('Tabs', () => {
  it('wraps Base UI tabs behind local styled primitives', () => {
    const markup = renderToStaticMarkup(
      <Tabs defaultValue="backend">
        <TabsList aria-label="Stack">
          <TabsTrigger value="backend">Backend</TabsTrigger>
        </TabsList>
        <TabsContent value="backend">Go API</TabsContent>
      </Tabs>,
    )

    expect(markup).toContain('Backend')
    expect(markup).toContain('Go API')
    expect(markup).toContain('data-orientation')
  })

  it('switches the selected panel when a tab is clicked', async () => {
    vi.stubGlobal('IS_REACT_ACT_ENVIRONMENT', true)
    const container = document.createElement('div')
    document.body.append(container)
    const root = createRoot(container)

    try {
      await act(async () => {
        root.render(
          <Tabs defaultValue="backend">
            <TabsList aria-label="Stack">
              <TabsTrigger value="backend">Backend</TabsTrigger>
              <TabsTrigger value="frontend">Frontend</TabsTrigger>
            </TabsList>
            <TabsContent value="backend">Go API</TabsContent>
            <TabsContent value="frontend">React UI</TabsContent>
          </Tabs>,
        )
      })

      const tabs = container.querySelectorAll<HTMLButtonElement>('[role="tab"]')
      expect(tabs[0].getAttribute('aria-selected')).toBe('true')

      await act(async () => tabs[1].click())

      expect(tabs[0].getAttribute('aria-selected')).toBe('false')
      expect(tabs[1].getAttribute('aria-selected')).toBe('true')
      const panelId = tabs[1].getAttribute('aria-controls')
      expect(panelId).toBeTruthy()
      const panel = document.getElementById(panelId!)
      expect(panel?.textContent).toBe('React UI')
      expect(panel?.hidden).toBe(false)
    } finally {
      await act(async () => root.unmount())
      container.remove()
      vi.unstubAllGlobals()
    }
  })
})
