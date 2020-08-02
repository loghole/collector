<template>
  <div class="columns p-2 pt-4">
    <div class="column page-menu">
      <!-- date -->
      <b-field>
        <DateTime
          :startTime="form.startTime"
          :end-time="form.endTime"
          :interval="form.interval"
          v-on:setStartTime="setStartTime"
          v-on:setEndTime="setEndTime"
          v-on:setInterval="setInterval"
        ></DateTime>
      </b-field>
      <!-- // date -->

      <!-- menu default params -->
      <b-field>
        <Menu
          v-on:setFormField="setFormField"
          v-on:setOperatorField="setOperatorField"
          :form="form"
          :operator="operator"
          :params="defaultParams"
        >
        </Menu>
      </b-field>
      <!-- // menu default params -->

      <!-- params -->
      <JSONValue
        v-for="(param, i) in params"
        :key="`param_${i}`"
        :param="param"
        :index="i"
        v-on:setJSONField="setJSONField"
        v-on:setJSONOperator="setJSONOperator"
        v-on:deleteJSONParam="deleteJSONParam"
      ></JSONValue>
      <!-- // params -->

      <!-- menu additional param -->
      <template v-if="showAdditionalParam">
        <b-field>
          <Menu
            v-on:setFormField="setFormField"
            v-on:setOperatorField="setOperatorField"
            :form="form"
            :operator="operator"
            :params="additionalParams"
          ></Menu>
        </b-field>
      </template>
      <!-- // menu additional param -->

      <!-- add new param  -->
      <b-field v-if="showAddParam">
        <b-input
          v-model="newParamName"
          placeholder="new param name"
          type="text"
          @keydown.native.enter="saveParam"
        >
        </b-input>
        <p class="control">
          <button class="button is-primary" @click="saveParam">Add</button>
        </p>
      </b-field>
      <!-- // add new param  -->

      <div class="buttons is-centered">
        <button
          class="button is-small is-outlined"
          @click="showAdditionalParam = !showAdditionalParam"
        >
          <b-icon
            :icon="showAdditionalParam ? 'eye-off' : 'eye'"
            size="is-small"
          >
          </b-icon>
          <span>other</span>
        </button>

        <!-- add param -->
        <button class="button is-small is-outlined" @click="newParam">
          <b-icon icon="plus" size="is-small"> </b-icon>
          <span>param</span>
        </button>
        <!-- // add param -->
      </div>
      <b-button :loading="loadingState" class="button is-primary is-fullwidth" @click="search(0)"
        >Search</b-button
      >
    </div>

    <div class="column">
      <div class="columns">
        <!-- Showed tags -->
        <div class="column">
          <b-taginput
            v-model="showedTags"
            :data="filteredTags"
            autocomplete
            :allow-new="true"
            :open-on-focus="true"
            placeholder="Showed tags"
            icon="label"
            @typing="setFilteredTag"
          ></b-taginput>
        </div>
        <!-- // Showed tags -->

        <div class="column">
          <b-field label="Search" label-position="on-border">
            <b-input
              placeholder="Search..."
              type="search"
              icon="magnify"
              icon-clickable
              class="w100"
              v-model="form.message"
            ></b-input>
            <p class="control">
              <b-button :loading="loadingState" class="button is-primary" @click="search(0)">Search</b-button>
            </p>
          </b-field>
        </div>
      </div>

      <!-- messages table -->
      <MessagesTable
        :activeTags="showedTags"
        :messages="messages"
      ></MessagesTable>
      <!-- // entry table -->
    </div>
  </div>
</template>

<script lang="ts">
import Vue from 'vue';
import DateTime from '@/components/DateTime.vue';
import Menu from '@/components/menu/Menu.vue';
import MessagesTable from '@/components/messages/MessagesTable.vue';
import JSONValue from '@/components/menu/JSONValue.vue';
import {
  Param, Form, ParamValue, SearchParam,
} from '@/types/view';

import { SingleParam, IntervalRegexp } from '@/const/const';
import FilterTags from '@/plugins/filter';

const footerHeight = 250;

export default Vue.extend({
  components: {
    DateTime,
    Menu,
    MessagesTable,
    JSONValue,
  },
  data() {
    return {
      operator: {
        level: '=' as string,
        namespace: '=' as string,
        source: '=' as string,
        traceID: '=' as string,
        host: '=' as string,
        buildCommit: '=' as string,
        configHash: '=' as string,
        remoteIP: '=' as string,
      } as Record<string, string>,
      form: {
        startTime: new Date(new Date().getTime() - 1000 * 15),
        endTime: null as Date | null,
        interval: '5m',
        namespace: [] as string[],
        source: [] as string[],
        traceID: [] as string[],
        host: [] as string[],
        level: [] as string[],
        buildCommit: [] as string[],
        configHash: [] as string[],
        remoteIP: [] as string[],
        message: '' as string,
      } as Form,
      defaultParams: [
        { key: 'level', name: 'Level', type: 'level' },
        { key: 'namespace', name: 'Namespace', type: 'namespace' },
        { key: 'source', name: 'Source', type: 'source' },
        { key: 'traceID', name: 'Trace ID' },
      ] as SearchParam[],
      additionalParams: [
        { key: 'host', name: 'Host', type: 'host' },
        { key: 'buildCommit', name: 'Build commit' },
        { key: 'configHash', name: 'Config hash' },
        { key: 'remoteIP', name: 'remote IP' },
      ] as SearchParam[],
      params: [] as Param[],
      showAddParam: false,
      newParamName: '',
      tagsInput: '',
      tags: [] as string[],
      showedTags: [],
      messages: [],
      showAdditionalParam: false,
      loadingState: false,
      loadLock: false,
    };
  },
  computed: {
    filteredTags(): string[] {
      return FilterTags(this.tags, this.showedTags, this.tagsInput);
    },
    scrollHeight() {
      return window.innerHeight + footerHeight;
    },
  },
  methods: {
    saveParam(): void {
      this.showAddParam = false;

      if (this.params.findIndex((p) => p.key === this.newParamName) !== -1) {
        // TODO higlight this param // 5 sec
        return;
      }

      this.params.push({
        type: 'json',
        key: this.newParamName,
        value: { item: '', list: [] as string[] },
        operator: '=',
      });
    },
    newParam() {
      this.showAddParam = true;
      this.newParamName = '';
    },
    setStartTime(val: Date): void {
      this.form.startTime = val;
    },
    setEndTime(val: Date): void {
      this.form.endTime = val;
    },
    setInterval(val: string): void {
      this.form.interval = val;
    },
    setFormField(key: string, val: string[]): void {
      this.form[key] = val;
    },
    setOperatorField(key: string, val: string): void {
      this.operator[key] = val;
    },
    setJSONField(idx: number, val: string | string[]): void {
      if (typeof val === 'string') {
        this.params[idx].value.item = val;
        return;
      }

      this.params[idx].value.list = val;
    },
    setJSONOperator(idx: number, val: string): void {
      if (
        !SingleParam.includes(this.params[idx].operator)
        === SingleParam.includes(val)
      ) {
        this.params[idx].value.list = [];
        this.params[idx].value.item = '';
      }

      this.params[idx].operator = val;
    },
    deleteJSONParam(idx: number): void {
      this.$nextTick(() => {
        this.params = this.params.filter((v, i) => i !== idx);
      });
    },
    setFilteredTag(text: string) {
      this.tagsInput = text;
    },
    convertInterval(val: string): Date {
      const values = IntervalRegexp.exec(val);

      let t = '';
      let num = 0;

      if (values !== null) {
        num = parseInt(values[1], 10);
        t = values[2] as string;
      }

      let offset = 0;

      switch (t) {
        case 'm':
        case 'min':
          offset = num * 60;
          break;

        case 'h':
        case 'hr':
        case 'hour':
          offset = num * 3600;
          break;

        case 'd':
        case 'day':
          offset = num * 3600 * 24;
          break;

        default:
          offset = num;
      }

      offset *= 1000;

      return new Date(new Date().getTime() - offset);
    },
    search(offset = 0, callback: any = null): void {
      let time = this.form.startTime as Date;

      if (this.form.interval !== '0') {
        time = this.convertInterval(this.form.interval as string);
      }

      const params = [
        {
          type: 'column',
          key: 'time',
          operator: '>=',
          value: {
            item: parseInt(
              (time.getTime() / 1000).toString(),
              10,
            ).toString(),
          } as ParamValue,
        },
      ] as Param[];

      if (this.form.endTime !== null) {
        params.push({
          type: 'column',
          key: 'time',
          operator: '<=',
          value: { item: this.form.endTime } as ParamValue,
        });
      }

      if (this.form.endTime !== null) {
        params.push({
          type: 'column',
          key: 'time',
          operator: '<=',
          value: { item: this.form.endTime } as ParamValue,
        });
      }

      if (this.form.message !== '') {
        params.push({
          type: 'column',
          key: 'message',
          operator: 'LIKE',
          value: { item: this.form.message } as ParamValue,
        });
      }

      [
        { mapKey: 'namespace', key: 'namespace', value: this.form.namespace },
        { mapKey: 'level', key: 'level', value: this.form.level },
        { mapKey: 'source', key: 'source', value: this.form.source },
        { mapKey: 'trace_id', key: 'trace_id', value: this.form.traceID },
        { mapKey: 'host', key: 'host', value: this.form.host },
        {
          mapKey: 'buildCommit',
          key: 'build_commit',
          value: this.form.buildCommit,
        },
        {
          mapKey: 'configHash',
          key: 'config_hash',
          value: this.form.configHash,
        },
        {
          mapKey: 'remoteIP',
          key: 'remote_ip',
          value: this.form.remoteIP,
        },
      ].forEach((h) => {
        const val = h.value as string[];
        if (val.length > 0) {
          params.push({
            type: 'column',
            key: h.key,
            operator: this.operator[h.mapKey] || '=',
            value: { list: h.value } as ParamValue,
          });
        }
      });

      this.params.forEach((param) => {
        params.push({
          type: 'key',
          key: param.key,
          operator: param.operator,
          value: param.value,
        });
      });

      // console.log(JSON.stringify(params));

      if (decodeURIComponent(window.location.hash) !== `#/logs/view${decodeURIComponent(this.getFullURL())}`) {
        this.$router.push(this.getFullURL());
      }

      const startUT = new Date().getTime();
      this.loadingState = true;

      Vue.axios
        .post('/api/v1/entry/list', { params, limit: 50, offset })
        .then((response) => {
          if (offset !== 0) {
            this.messages = this.messages.concat(response.data.data);
          } else {
            this.messages = response.data.data;
          }

          this.setTags(response.data.data);
        })
        .catch((e) => {
          this.$buefy.notification.open({
            duration: 5000,
            message: e.message,
            position: 'is-bottom-right',
            type: 'is-danger',
            hasIcon: true,
          });

          console.error(e);
        }).then(() => {
          const diff = new Date().getTime() - startUT;

          if (callback !== null) {
            callback();
          }

          setTimeout(() => {
            this.loadingState = false;
          }, 250 - diff);
        });
    },
    getFullURL(): string {
      return `?form=${this.getURL(this.form)}&params=${this.getURL(this.params)}&tags=${this.getURL(this.showedTags)}`;
    },
    getURL(param: Form | Param[] | string[]): string {
      return encodeURIComponent(JSON.stringify(param));
    },
    setTags(list: Array<any>): void {
      const h = {} as Record<string, boolean>;

      list.forEach((l: Record<string, number>) => {
        Object.keys(l).forEach((k) => {
          h[k] = true;
        });
      });

      this.tags = Object.keys(h);
    },
    handleScroll() {
      if (
        this.loadLock === false
        && document.body.scrollHeight - this.scrollHeight - window.scrollY < 0
      ) {
        this.loadLock = true;

        this.search(this.messages.length, () => {
          this.loadLock = false;
        });
      }
    },
  },
  created() {
    window.addEventListener('scroll', this.handleScroll);

    if (this.$route.query.params) {
      this.params = JSON.parse(decodeURIComponent(this.$route.query.params as string));
    }

    if (this.$route.query.tags) {
      this.showedTags = JSON.parse(decodeURIComponent(this.$route.query.tags as string));
    }

    if (this.$route.query.form) {
      const form = JSON.parse(decodeURIComponent(this.$route.query.form as string));
      form.startTime = new Date(form.startTime);

      this.form = form;
      this.search();

      return;
    }

    Vue.axios
      .post('/api/v1/entry/list', { limit: 100 })
      .then((response) => {
        this.messages = response.data.data;

        this.setTags(response.data.data);

        console.log(response.data);
      })
      .catch((e) => {
        console.error(e);
      });
  },
  beforeDestroy() {
    window.removeEventListener('scroll', this.handleScroll);
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

.w100 {
  width: 100%;
}

.page {
  &-menu {
    max-width: 210px;
    min-width: 150px;
  }
}

.is-relative {
  position: relative;
}
</style>
