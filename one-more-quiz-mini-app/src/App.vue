<template>
  <div v-if='currentScreen === "quiz"'>
    <QuizScreen v-on:to-statistic='currentScreen = "statistic"' />
  </div>
  <div v-else>
    <StatisticScreen v-on:to-quiz='currentScreen = "quiz"' />
  </div>
</template>

<script>

/* Import a component for each screen, and a property 'currentScreen' for toggling screens by id */

import QuizScreen from './components/QuizScreen.vue'
import StatisticScreen from './components/StatisticScreen.vue'

var screens = ["quiz", "statistic"];

export default {
  name: 'App',
  data() {
    return {
      screens: screens,
      currentScreen: screens[0],
    };
  },
  components: {
    QuizScreen,
    StatisticScreen
  }
}
</script>

<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  margin-top: 32px;
  /* Prevent webview from selecting and resizing on mobile devices */
  -webkit-user-select: none; /* Safari */
  -ms-user-select: none; /* IE 10 and IE 11 */
  user-select: none; /* Standard syntax */
  -webkit-touch-callout: none; /* prevent callout to copy image, etc when tap to hold */
  -webkit-text-size-adjust: none; /* prevent webkit from resizing text to fit */
}

</style>
