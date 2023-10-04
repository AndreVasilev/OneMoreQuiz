<template>
  <div class="quiz">

    <div class="header">
      <label class="profile-button" @click="toStatistic">{{ user.Score }}</label>
    </div>
    
    <div class="question-container">
      <div class="timer-container">
        <h1 class="timer">{{ timer > 0 ? timer : '' }}</h1>
      </div>
      <h3 class="question">{{ current.question }}</h3>
    </div>
    
    <div class="answers-container">
      <template v-for="answer in current.answers" :key="answer.id">
        <label v-bind:class="[answer.id === checked 
            ? (answer.id === current.answer 
                ? 'answer-true' 
                : 'answer-false') 
            : (stopped && answer.id === current.answer 
                ? 'answer-true' 
                : 'answer')]" @click="checkAnswer(answer.id)">
          {{answer.text}}
        </label>
      </template>
    </div>
    
    <!-- Uncomment when running in develop mode locally to enable MainButton functionality -->
    <!-- <button v-on:click="startOneMore">One more!</button> -->

  </div>
</template>

<script>
import axios from 'axios';
import AppConfig from './AppConfig.vue'

export default {
  name: "QuizScreen",
  data() {
    return {
      /* List of received questions */
      questions: null,
      /* Selected question that is shown to user*/
      current: {
        id: null,
        question: null,
        answers: [
            {id: "A", text: null},
            {id: "B", text: null},
            {id: "C", text: null},
            {id: "D", text: null}
        ],
        answer: null
      },
      /* User model for presenting actual Score */
      user: {
        Id: null,
        LastQuestionId: null,
        SuccessAnswers: null,
        Score: null,
      },
      /* Flag: user did select an answer to current question */
      checked: '',
      /* Timer: its value is set in AppConfig.vue */
      timer: 0,
      timerInterval: null,
      /* Flag: timer expired, there is no time to answer current question */
      stopped: false
    };
  },
  methods: {
    /* Open Statistic screen */
    toStatistic: function() {
      this.mainButtonSetup(false);
      this.$emit('to-statistic')
    },
    /* Setup MainButton */
    mainButtonSetup: function(isVisible) {
      window.Telegram.WebApp.MainButton.isVisible = isVisible;
      if (isVisible) {
        window.Telegram.WebApp.MainButton.setText("One more!");
        var self = this;
        window.Telegram.WebApp.MainButton.onClick(function() {
          self.startOneMore();
        });
      }
    },
    /* Tell Telegram to play Haptic Feedback if available */
    hapticFeedback: function() {
      if (window.Telegram.WebApp.isVersionAtLeast("6.1")) {
        window.Telegram.WebApp.HapticFeedback.selectionChanged()
      }
    },
    /* Get questions for ther user from backend */
    getQuestions: function() {
      var self = this;
      axios
        .post(AppConfig.api.question, JSON.stringify({
          "initData": AppConfig.tgInitData()
        }))
        .then(function(response) {
          self.questions = response.data;
          self.startOneMore();
        })
    },
    /* Post to the backend current question id */
    postUser: function(questionId) {
      var self = this;
      var url = AppConfig.api.user + `?question_id=${questionId}`
      axios
        .post(url, JSON.stringify({
          "initData": AppConfig.tgInitData()
        }))
        .then(function(response) {
          self.user = response.data;
        })
    },
    /* Post to the backend user's answer for current question */
    postUserAnswer: function(questionId, succeed, score) {
      var self = this;
      var url = AppConfig.api.user + `?question_id=${questionId}&succeed=${succeed}&score=${score}`
      axios
        .post(url, JSON.stringify({
          "initData": AppConfig.tgInitData()
        }))
        .then(function(response) {
          self.user = response.data;
        })
    },
    /* Check answer when user did select one */
    checkAnswer: function(checked) {
      if (this.stopped || this.checked != '') {
        return
      }
      this.hapticFeedback();
      this.checked = checked;
      var self = this;
      var t = setInterval(function() {
        clearInterval(t)
        self.stop()
      }, 300)
    },
    /* Show next question, reset timer and flags */
    startOneMore: function() {
      if (this.questions == null || this.questions.length < 1) {
        return;
      }
      var next = this.questions.shift();
      this.current = {
        id: next.Id,
        question: next.Question,
        answers: [
            {id: "A", text: next.A},
            {id: "B", text: next.B},
            {id: "C", text: next.C},
            {id: "D", text: next.D}
        ],
        answer: next.Answer
      };
      this.hapticFeedback();
      this.postUser(next.Id);
      this.checked = ''
      this.timer = AppConfig.timer
      this.stopped = false
      var self = this;
      clearInterval(this.timerInterval)
      this.timerInterval = setInterval(function() {
        if (self.stopped) {
          return
        }
        self.timer -= 1
        if (self.timer <= 0) {
          clearInterval(self.timerInterval)
          self.stop()
        }
      }, 1000)
    },
    /* Stop timer and post user's answer to the backend when a user did select an answer or the timer expired */
    stop: function() {
      this.stopped = true
      clearInterval(this.timerInterval)
      var succeed = this.checked == this.current.answer
      this.postUserAnswer(this.current.id, succeed, succeed ? this.timer : 0);
      this.timer = 0
      this.mainButtonSetup(true);
    }
  },
  /* Setup and reload initial data when screen created */
  created() {
    this.mainButtonSetup(false);
    if (this.questions == null) {
      this.getQuestions();
    }
  },
};
</script>

<style>
  .header {
    text-align: right;
  }

  .profile-button {
    height: 44px;
    border-radius: 12px;
    color: var(--tg-theme-link-color);
    text-align: center;
    cursor: pointer;
    box-shadow: 0px 2px 4px 2px rgba(0, 0, 0, 0.1);
    font-weight: 500;
    background-color: var(--tg-theme-bg-color);
    padding: 8px 16px 8px 16px;
    margin-right: 16px;
    line-height: 44px;
  }

  .question-container {
    position: relative;
    background-color: var(--tg-theme-bg-color);
    border-radius: 16px;
    box-shadow: 0px 8px 16px 8px rgba(0, 0, 0, 0.1);
    padding: 8px;
    margin: 64px 16px 0px 16px;
  }

  .timer-container {
    position: absolute;
    top: -48px;
    left: calc(50% - 48px);
    width: 96px;
    height: 96px;
    background-color: var(--tg-theme-bg-color);
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    box-shadow: 0px 4px 8px 4px rgba(0, 0, 0, 0.05);
  }

  .timer {
    color: var(--tg-theme-text-color);
  }

  .question {
    margin-top: 64px;
    color: var(--tg-theme-text-color);
  }

  .answers-container {
    display: flex;
    gap: 8px; /* Adjust the gap between buttons as needed */
    flex-direction: column;
    margin: 20px 16px 0px 16px;
  }

  .answers-container label {
    display: inline-block;
    width: 100%;
    height: 44px;
    border-radius: 12px;
    color: var(--tg-theme-text-color);
    text-align: center;
    line-height: 44px;
    cursor: pointer;
    box-shadow: 0px 2px 4px 2px rgba(0, 0, 0, 0.1);
    font-weight: 500;
  }

  .answer {
    background-color: var(--tg-theme-bg-color);
  }
  .answer:active, .profile-button:active {
    opacity: 0.7;
  }

  .answer-true {
    background-color: #00A47899;
  }

  .answer-false {
    background-color: #E4484899;
  }
</style>
