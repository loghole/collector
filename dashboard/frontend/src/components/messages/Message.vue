<template>
  <tbody>
    <tr>
      <td>
        <button aria-expanded="true" aria-label="Toggle row details" class="btn" @click="showJSON">
          <b-icon :icon="jsonBlockIsShowed ? 'menu-down' : 'menu-down-outline'"></b-icon>
        </button>
      </td>
      <td style="width: 166px">{{ buildDatetime(message.time) }}</td>

      <!-- tags with values -->
      <template v-if="activeTags.length > 0">
        <td v-for="(tag, i) in activeTags" :key="`tag_${i}`">
          {{ prepareText(tag, message[tag]) }}
        </td>
      </template>
      <td v-else>
        <dl class="tag-source">
          <template v-for="(value, name, i) in message">
            <template v-if="showField(name, value)">
            <dt
              class="tag-name"
              :key="`tag_name_${i}`">{{ name }}:</dt>
            <dd
              class="tag-value"
              :key="`tag_value_${i}`">{{ prepareText(name, value) }}</dd>
            </template>
          </template>
        </dl>
      </td>
    </tr>
    <!-- // tags with values -->

    <!-- json block -->
    <tr v-if="jsonBlockIsShowed">
      <td :colspan="activeTags.length + 3">
        <JSONCode :data="message.params"></JSONCode>
      </td>
    </tr>
    <!-- // all json block -->
  </tbody>
</template>

<script lang="ts">
import Vue from 'vue';
import JSONCode from '@/components/messages/JSONCode.vue';

export default Vue.extend({
  components: {
    JSONCode,
  },
  props: {
    message: {
      type: Object,
      required: true,
    },
    activeTags: {
      type: Array,
      required: true,
    },
  },
  data() {
    return {
      jsonBlockIsShowed: false,
    };
  },
  methods: {
    prepareText(tag: string, text: string): string {
      switch (tag) {
        case 'time':
          return this.buildDatetime(text);
        case 'level':
          return text.toUpperCase();
        default:
          return text;
      }
    },
    buildDatetime(text: string): string {
      return new Date(text).toLocaleString(window.navigator.language);
    },
    showField(name: string, value: string): boolean {
      return name !== 'params' && value !== '';
    },
    showJSON(): void {
      this.jsonBlockIsShowed = !this.jsonBlockIsShowed;
    },
  },
});
</script>

<style lang="scss" scoped>
  /* TODO: fix it */
  .btn {
    background: none;
    border: none;
    padding: 0;
    margin: 0;
    outline: none;
    font-size: inherit;
    color: inherit;
    border-radius: 0;
    :active, :focus {
      background: none;
      outline: none;
      box-shadow: none;
      border-color: #dbdbdb;
    }
  }
  .tag-source {
    margin-bottom: 0;
    line-height: 2em;
    word-break: break-word;
    padding: 8px;
    overflow: hidden;
    max-height: 100px !important;
    display: inline-block;
  }
  .tag-name {
    background-color: rgba(3, 130, 217, 0.14);
    padding: 2px 4px;
    margin-right: 4px;
    word-break: normal;
    border-radius: 4px;
    display: inline;
  }
  .tag-value {
    display: inline-flex;
    word-break: break-word;
    margin-right: 5px;
  }
</style>
