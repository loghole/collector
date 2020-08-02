<template>
    <b-sidebar
      type="is-light"
      :fullheight="true"
      :overlay="false"
      :open.sync="open"
    >
      <div class="p-2 pt-4">
        <b-field label="Operator" label-position="on-border">
          <b-autocomplete
            v-model="form.operator"
            placeholder="e.g. >="
            :data="filteredOperators"
            :open-on-focus="true"
            @select="option => (selected = option)"
          >
          </b-autocomplete>
        </b-field>

        <b-field label="Key" label-position="on-border">
          <input
            class="input"
            v-model="form.key"
            type="text"
            placeholder="Field name"
          />
        </b-field>

        <b-field
          label="Value"
          label-position="on-border"
          v-if="isListValue(form.operator)"
        >
          <b-taginput
            v-model="form.value.list"
            autocomplete
            :allow-new="true"
            placeholder="Value"
            icon="label"
          >
          </b-taginput>
        </b-field>
        <b-field label="Value" label-position="on-border" v-else>
          <input
            class="input"
            v-model="form.value.item"
            type="text"
            placeholder="Value"
          />
        </b-field>

        <button
          class="button is-small is-fullwidth is-outlined is-success"
          @click="save()"
        >
          Add
        </button>
      </div>
    </b-sidebar>
</template>

<script lang="ts">
import Vue from 'vue';
import { Param } from '../types/view';

export default Vue.extend({
  props: {
    value: {
      type: Boolean,
      required: true,
    },
  },
  data() {
    return {
      form: {
        operator: '',
        type: '',
        key: '',
        value: {
          item: '',
          list: [],
        },
      } as Param,
      operators: ['=', '!=', '>', '<', '>=', '<=', 'LIKE', 'NOT LIKE'],
    };
  },
  computed: {
    open: {
      get() {
        return this.value;
      },
      set(newValue) {
        this.$emit('input', newValue);
      },
    },
    filteredOperators(): string[] {
      return this.operators.filter(
        (option: string) => option.toLowerCase().indexOf(this.form.operator.toLowerCase()) >= 0,
      );
    },
  },
  methods: {
    isListValue(operator: string): boolean {
      return ['=', '!=', 'LIKE', 'NOT LIKE'].includes(operator);
    },
    save() {
      this.$emit('save', this.form);

      this.form = {
        operator: '',
        type: '',
        key: '',
        value: { item: '', list: [] },
      };
    },
  },
});
</script>

<style lang="scss" scoped>
.p-2 {
  padding: 0.5rem;
}
.pt-4 {
  padding-top: 1rem;
}
</style>
