<template>
  <div v-if="hasContent" class="table-container" >
    <table class="table is-striped is-narrow is-hoverable is-fullwidth">
      <thead>
      <tr>
        <th style="width: 24px"></th>
        <th>Time</th>
        <template v-if="activeTags.length > 0">
          <th v-for="(tag, i) in activeTags" :key="i">
            {{ tag }}
          </th>
        </template>
        <template v-else>
          <th>_source</th>
        </template>
      </tr>
      </thead>
      <Message
        v-for="(message, i) in messages"
        :active-tags="activeTags"
        :message="message"
        :key="i">
      </Message>
    </table>
  </div>
  <section v-else class="section">
    <div class="content has-text-grey has-text-centered">
      <p>
        <b-icon
          icon="emoticon-sad"
          size="is-large">
        </b-icon>
      </p>
      <p>Nothing here.</p>
    </div>
  </section>
</template>

<script lang="ts">
import Vue from 'vue';
import Message from '@/components/messages/Message.vue';

export default Vue.extend({
  components: {
    Message,
  },
  props: {
    activeTags: {
      type: Array,
      required: true,
    },
    messages: {
      type: Array,
      required: true,
    },
  },
  computed: {
    hasContent(): boolean {
      return this.messages.length > 0;
    },
  },
});
</script>
