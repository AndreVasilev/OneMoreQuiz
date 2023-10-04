<template>
  <div class="statistic">

    <h1 class="main-title">Statistic</h1>
    
    <div class="container">
      <div class="row">
        <label class="title">Rating</label>
        <label class="value">{{ user.Score }}</label>
      </div>
      <div class="row">
        <label class="title">Rating position</label>
        <label class="value">{{ user.Position }}</label>
      </div>
      <div class="row">
        <label class="title">Total questions</label>
        <label class="value">{{ user.LastQuestionId }}</label>
      </div>
      <div class="row">
        <label class="title">Right answers</label>
        <label class="value">{{ user.SuccessAnswers }}</label>
      </div>
    </div>

    <!-- Uncomment when running in develop mode locally to enable BackButton functionality -->
    <!-- <button @click="$emit('to-quiz')">Back</button> -->

  </div>
</template>

<script>
import axios from 'axios';
import AppConfig from './AppConfig.vue'

export default {
  name: "StatisticScreen",
  data() {
    return {
      /* Users's model */
      user: {
        Id: null,
        LastQuestionId: null,
        SuccessAnswers: null,
        Score: null,
        Position: null
      },
    };
  },
  methods: {
    /* Reload data */
    reloadData() {
      axios
        .post(AppConfig.api.user, JSON.stringify({
          "initData": AppConfig.tgInitData()
        }))
        .then(response => (this.user = response.data));
    },
    /* Setup BackButton */
    backButtonSetup: function() {
      var self = this;
      window.Telegram.WebApp.BackButton.isVisible = true;
      window.Telegram.WebApp.BackButton.onClick(function() {
        window.Telegram.WebApp.BackButton.isVisible = false;
        self.$emit('to-quiz')
      });
    },
  },
  /* Setup and reload initial data when screen mounted */
  mounted() {
    this.reloadData();
    this.backButtonSetup();
  }
};
</script>

<style>
  .main-title {
    color: var(--tg-theme-text-color);
  }

  .container {
    position: relative;
    background-color: var(--tg-theme-bg-color);
    border-radius: 16px;
    box-shadow: 0px 8px 16px 8px rgba(0, 0, 0, 0.1);
    padding: 16px;
    margin: 16px 16px 0px 16px;
    font-weight: 500;
  }

  .row {
    display: flex;
    justify-content: space-between;
  }

  .title {
    margin-top: 12px;
    color: var(--tg-theme-hint-color);
  }

  .value {
    margin-top: 12px;
    color: var(--tg-theme-text-color);
  }
</style>