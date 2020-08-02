<template>
  <div>
    <b-field
      label-position="on-border"
      class="is-relative"
      v-for="(item, i) in params"
      :key="`param${i}`"
    >
      <TagInput
        :value="form[item.key]"
        @input="setField(item.key, $event)"
        :placeholder="item.name"
        :type="item.type"
      >
      </TagInput>

      <LableList
        :isMultiple="true"
        :value="operator[item.key]"
        @input="setOperator(item.key, $event)"
        :name="item.name"
      ></LableList>
    </b-field>
  </div>
</template>

<script lang="ts">
import Vue from 'vue';
import TagInput from '@/components/TagInput.vue';
import LableList from '@/components/LabelList.vue';

export default Vue.extend({
  components: {
    TagInput,
    LableList,
  },
  props: {
    params: {
      type: Array,
      required: true,
    },
    form: {
      type: Object,
      required: true,
    },
    operator: {
      type: Object,
      required: true,
    },
  },
  methods: {
    setField(name: string, value: string[]) {
      this.$emit('setFormField', name, value);
    },
    setOperator(name: string, value: string) {
      this.$emit('setOperatorField', name, value);
    },
  },
});
</script>
