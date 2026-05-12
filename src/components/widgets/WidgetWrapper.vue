<template>
  <div class="widget-wrapper" :class="[`widget-wrapper--${variant}`, { 'widget-wrapper--no-header': !headerVisible }]">
    <div v-if="headerVisible" class="widget-header">
      <div class="widget-header__leading">
        <span v-if="icon" class="widget-header__icon" aria-hidden="true">{{ icon }}</span>
        <span v-if="title" class="widget-header__title">{{ title }}</span>
      </div>

      <div class="widget-header__trailing">
        <slot name="actions" />
        <span
          v-if="draggable"
          class="widget-drag-handle"
          :title="dragHandleTitle"
          aria-label="Drag to reposition widget"
          role="button"
        >⠿</span>
      </div>
    </div>

    <div class="widget-body" :style="bodyStyle">
      <slot />
    </div>
  </div>
</template>

<script setup>
import { computed, useSlots } from 'vue'

const slots = useSlots()

const props = defineProps({
  /** Optional display title shown in the header */
  title: {
    type: String,
    default: '',
  },

  /** Optional icon/emoji shown before the title */
  icon: {
    type: String,
    default: '',
  },

  /** Controls whether the header bar is rendered at all. When true (default), header auto-hides if it would be empty. */
  showHeader: {
    type: Boolean,
    default: true,
  },

  /** Controls whether the drag handle is shown */
  draggable: {
    type: Boolean,
    default: true,
  },

  /** Accessible label for the drag handle */
  dragHandleTitle: {
    type: String,
    default: 'Drag to reposition',
  },

  /**
   * Visual variant.
   * 'default' | 'flat' | 'transparent'
   */
  variant: {
    type: String,
    default: 'default',
    validator: (v) => ['default', 'flat', 'transparent'].includes(v),
  },

  /** Optional padding inside the widget body. Defaults to 0. */
  bodyPadding: {
    type: String,
    default: '0',
  },
})

const headerVisible = computed(() => {
  if (!props.showHeader) return false
  return !!(props.title || props.icon || slots.actions || props.draggable)
})

const bodyStyle = computed(() => ({
  padding: props.bodyPadding,
}))
</script>

<style scoped>
/* ── Layout ──────────────────────────────────────────────────────────────── */

.widget-wrapper {
  display: flex;
  flex-direction: column;
  width: 100%;
  height: 100%;
  overflow: hidden;
  box-sizing: border-box;

  /* Theming tokens — override these from the outside as needed */
  --widget-bg: #1a1a1a;
  --widget-header-bg: #202024;
  --widget-border-color: #2a2a2e;
  --widget-handle-color: #555;
  --widget-title-color: #aaa;
  --widget-radius: 8px;

  background-color: var(--widget-bg);
  border-radius: var(--widget-radius);
}

/* ── Variants ─────────────────────────────────────────────────────────────── */

.widget-wrapper--flat {
  --widget-bg: transparent;
  --widget-header-bg: transparent;
  --widget-border-color: transparent;
  border: 1px solid var(--widget-border-color);
}

.widget-wrapper--transparent {
  --widget-bg: transparent;
  --widget-header-bg: transparent;
  --widget-border-color: transparent;
}

/* ── Header ───────────────────────────────────────────────────────────────── */

.widget-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 4px 8px;
  flex-shrink: 0;
  background-color: var(--widget-header-bg);
  border-bottom: 1px solid var(--widget-border-color);
  min-height: 28px;
}

.widget-header__leading {
  display: flex;
  align-items: center;
  gap: 6px;
  overflow: hidden;
}

.widget-header__icon {
  font-size: 14px;
  line-height: 1;
  flex-shrink: 0;
}

.widget-header__title {
  font-size: 12px;
  font-weight: 500;
  color: var(--widget-title-color);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  user-select: none;
}

.widget-header__trailing {
  display: flex;
  align-items: center;
  gap: 4px;
  flex-shrink: 0;
  margin-left: 8px;
}

/* ── Drag handle ──────────────────────────────────────────────────────────── */

.widget-drag-handle {
  color: var(--widget-handle-color);
  font-size: 14px;
  letter-spacing: 2px;
  user-select: none;
  line-height: 1;
  cursor: grab;
  padding: 2px 4px;
  border-radius: 4px;
  transition: color 0.15s ease;
}

.widget-drag-handle:hover {
  color: #888;
}

.widget-drag-handle:active {
  cursor: grabbing;
}

/* ── Body ─────────────────────────────────────────────────────────────────── */

.widget-body {
  flex: 1;
  overflow: hidden;
  min-height: 0;
  box-sizing: border-box;
}

/*
 * Force direct slot children to fill the body.
 * Kept as a catch-all but prefer explicit sizing inside child components
 * when possible.
 */
.widget-body > :deep(*) {
  width: 100%;
  height: 100%;
  min-width: 0;
  box-sizing: border-box;
}

/* Reset margin injected by card libraries (PrimeVue, etc.) */
.widget-body > :deep(.card),
.widget-body > :deep(.p-card) {
  width: 100% !important;
  height: 100% !important;
  min-width: 0 !important;
  margin: 0 !important;
}

/* Chart containers fill remaining space */
.widget-body > :deep(.chart-container) {
  height: 100%;
  min-height: 0;
}
</style>