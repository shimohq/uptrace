<template>
  <span>
    <v-btn
      v-if="span.statusCode === 'error'"
      icon
      :title="`${xkey.spanStatusCode} = 'error'`"
      class="mr-1"
      :class="{ 'cursor-default': !clickable }"
      @click.stop="$emit('click:chip', { key: xkey.spanStatusCode, value: 'error' })"
    >
      <v-icon color="error"> mdi-alert-circle-outline </v-icon>
    </v-btn>

    <v-chip
      v-for="(chip, i) in chips"
      :key="chip.key"
      color="light-blue lighten-5"
      label
      small
      :class="{ 'ml-1': i > 0, 'mb-1': i > 0, 'cursor-default': !clickable }"
      :title="`${chip.key}: ${chip.value}`"
      @click.stop="$emit('click:chip', chip)"
    >
      {{ chip.text }}
    </v-chip>

    <v-chip v-if="events.length" color="blue lighten-5" label small class="ml-1">
      <strong class="mr-1">{{ events.length }}</strong>
      <span>{{ events.length === 1 ? 'event' : 'events' }}</span>
    </v-chip>
  </span>
</template>

<script lang="ts">
import { truncate } from 'lodash'
import { defineComponent, computed, PropType } from '@vue/composition-api'

// Utilities
import { xkey } from '@/models/otelattr'
import { AttrMap, Span } from '@/models/span'

export interface SpanChip {
  key: string
  value: any
  text: string
}

export default defineComponent({
  name: 'SpanChips',

  props: {
    span: {
      type: Object as PropType<Span>,
      required: true,
    },
    showOperation: {
      type: Boolean,
      default: false,
    },
    traceMode: {
      type: Boolean,
      default: false,
    },
    clickable: {
      type: Boolean,
      default: false,
    },
  },

  setup(props) {
    const events = computed((): Span[] => {
      return props.span?.events ?? []
    })

    const chips = computed(() => {
      if (props.traceMode) {
        return traceChips(props.span.attrs)
      }

      const chips: SpanChip[] = []

      const service = props.span.attrs[xkey.serviceName]
      if (service) {
        chips.push({ key: xkey.serviceName, value: service, text: service })
      }

      pushKindChip(chips, props.span.attrs)

      const file = props.span.attrs[xkey.codeFilepath]
      if (file) {
        chips.push({ key: xkey.codeFilepath, value: file, text: shortFile(file) })
      }

      pushHttpStatusChip(chips, props.span.attrs)

      return chips
    })

    return { xkey, events, chips }
  },
})

function traceChips(attrs: AttrMap): SpanChip[] {
  const chips: SpanChip[] = []

  pushSystemChip(chips, attrs)
  pushKindChip(chips, attrs)
  pushHttpStatusChip(chips, attrs)

  return chips
}

function pushSystemChip(chips: SpanChip[], attrs: AttrMap) {
  const spanSystem = attrs[xkey.spanSystem]
  if (spanSystem && spanSystem !== xkey.internalSystem) {
    chips.push({ key: xkey.spanSystem, value: spanSystem, text: spanSystem })
  }
}

function pushKindChip(chips: SpanChip[], attrs: AttrMap) {
  const kind = attrs[xkey.spanKind]
  if (kind && kind !== 'internal') {
    chips.push({ key: xkey.spanKind, value: kind, text: kind })
  }
}

function pushHttpStatusChip(chips: SpanChip[], attrs: AttrMap) {
  const httpCode = attrs[xkey.httpStatusCode]
  if (typeof httpCode === 'number' && httpCode != 0 && (httpCode < 200 || httpCode >= 300)) {
    chips.push({ key: xkey.httpStatusCode, value: httpCode, text: String(httpCode) })
  }
}

function shortFile(s: string): string {
  let ind = s.lastIndexOf('/')
  if (ind !== -1) {
    s = s.slice(ind + 1)
  }
  return truncate(s)
}
</script>

<style lang="scss" scoped>
.cursor-default {
  cursor: default;
}
</style>
