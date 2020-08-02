<template>
  <div class="label">
    {{ name }}
    <b-dropdown aria-role="list">
      <span slot="trigger" role="button" class="pointer has-text-success"
        >{{ value }}
      </span>

      <b-dropdown-item
        aria-role="listitem"
        v-for="(param, i) in params"
        :key="`item${i}`"
        @click="setOperator(param)"
        >{{ param }}</b-dropdown-item
      >
    </b-dropdown>
  </div>
</template>

<script lang="ts">
import Vue from 'vue';
import { OperatorsSingle, OperatorMultiple } from '@/const/const';

export default Vue.extend({
  props: {
    value: {
      type: String,
      required: true,
    },
    isMultiple: {
      type: Boolean,
      required: true,
    },
    name: {
      type: String,
      default: '',
    },
  },
  methods: {
    setOperator(value: string): void {
      this.$emit('input', value);
    },
  },
  computed: {
    params(): string[] {
      return this.isMultiple ? OperatorMultiple : OperatorsSingle;
    },
  },
});
</script>

<style lang="scss" scoped>
.pointer {
  cursor: pointer;
}

.label {
  position: absolute;
  top: -0.775em;
  left: 1em;
  font-size: calc(1rem * 3 / 4);
  background-color: white;
  padding-left: 0.125em;
  padding-right: 0.125em;

  color: #363636;
  font-weight: 700;
}
</style>
