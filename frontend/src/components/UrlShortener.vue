<template>
  <div class="main-container">
    <div class="url-form-container">
      <form id="url-form" @submit="validateUrlForm" action="" method="POST">
        <label class="url-form-label">Enter URL:<br/>
          <p v-if="errors.url_error" class="form-error">Error: {{ errors.url_error }}</p>
          <input v-model.trim="original_url" v-bind:class="{ 'field-error': errors.url_error }" class="url-input" size="80" placeholder="http://example.com">
        </label><br/>
        <label class="url-input">How many days URL should be kept? Keep empty to delete right after the first usage:<br/>
          <p v-if="errors.keep_days_error" class="form-error">Error: {{ errors.keep_days_error }}</p>
          <input v-model.number="url_keep_days" v-bind:class="{ 'field-error': errors.keep_days_error }" class="url-keep-input" type="number">
        </label><br/>
        <p v-if="errors.fetch_error" class="form-error">Error: {{ errors.fetch }}</p>
        <input type="submit" class="button" value="Get short URL code">
      </form>
    </div>

    <div class="url-code-container">
      <label class="url-code-label">Short URL:
        <input v-model="url_code" class="url-result" size="35" readonly>
      </label><br/>
      <label class="url-keep-until-label">Would be kept until:
        <input v-model="keep_until" class="url-result" readonly>
      </label><br/>
      <button @click="copyToClipboard" class="button">Copy URL to clipboard</button><br/>
    </div>
  </div>
</template>

<script>
import axios from 'axios'
import Clipboard from 'v-clipboard'
import Vue from 'vue'

Vue.use(Clipboard);

const BACKEND_HOST = process.env.VUE_APP_BACKEND_HOST;


export default {
  name: "UrlShortener",

  data() {
    return {
      original_url: null,
      url_keep_days: null,
      errors: {},
      url_code: null,
      keep_until: null
    };
  },

  methods: {
    getUrlCode() {
      this.errors.fetch_error = false;

      let data = {url: this.original_url};
      if (this.url_keep_days !== null || this.url_keep_days !== '') {
        data['keep_for_days'] = this.url_keep_days;
      }

      axios.post(`${BACKEND_HOST}/url/add`, data)
        .then((resp) => {
          if (!resp.data.hasOwnProperty('url')) {
            this.errors.fetch_error = true;
            this.errors['fetch'] = "Internal error.";

            return;
          }

          this.url_code = resp.data.url;
          this.keep_until = resp.data.valid_until.slice(0, resp.data.valid_until.indexOf('.'));
        })
        .catch((error) => {
          this.errors.fetch_error = true;
          this.errors['fetch'] = error.response.data.error;
        });
    },

    copyToClipboard() {
      this.$clipboard(this.url_code);
    },

    validateUrlForm(e) {
      e.preventDefault();
      this.errors = {form_error: false};

      if (this.original_url === null || this.original_url.length === 0) {
        this.errors.form_error = true;
        this.errors['url_error'] = 'This field may not be empty!'
      }

      if (this.url_keep_days !== null && this.url_keep_days < 0) {
        this.errors.form_error = true;
        this.errors['keep_days_error'] = 'This field must contain only positive numbers!'
      }

      if (this.errors.form_error) {
        return false;
      }

      this.getUrlCode();

      return true;
    }
  }
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
.url-input,
.url-keep-input,
.url-result {
  margin: 10px 0;
  font-size: 18px;
  border-radius: 2px;
  border-width: 1px;
  border-color: #e0e0e0;
}

.url-result {
  border: none;
  background-color: #f5f5f5;
  margin: 10px 10px;
  padding: 3px;
}

.form-error {
  color: #820b0b;
  margin: 0;
}

.field-error {
  border-color: #d51616;
}

.button {
  color: #d5eeff;
  background-color: #4297c9;
  border-radius: 12px;
  border: none;
  font-size: 16px;
  font-weight: bold;
  padding: 10px 20px;
  width: 235px;
  margin: 10px 0;
}

.button:hover {
  color: #405a6b;
  background-color: #61b2e0;
}
.button:active {
  color: #566874;
  background-color: #a6daf8;
}
</style>
