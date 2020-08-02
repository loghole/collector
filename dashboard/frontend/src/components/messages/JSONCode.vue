<template>
  <pre class="code-block">
    <!-- TODO: span with color for strings and other elements -->
    <code class="json">{{ buildJSON() }}</code>
  </pre>
</template>

<script lang="ts">
import Vue from 'vue';

export default Vue.extend({
  props: {
    data: {
      type: Object,
      required: true,
    },
  },
  methods: {
    buildJSON(): string {
      return JSON.stringify(this.data, undefined, 4);
    },
  },
});
</script>

<style lang="scss" scoped>
  .code-block {
    white-space: inherit;
    padding: 8px;
  }
  .json {
    white-space: pre;
    font-family: "Roboto Mono", Consolas, Menlo, Courier, monospace;
    letter-spacing: normal;
    display: block;
    line-height: 1.5;
    font-weight: 400;
    font-size: inherit;
  }
</style>
