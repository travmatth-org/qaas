import Vue from 'vue';
import Buefy from 'buefy';

// import API from '../api';
// eslint-disable-next-line import/extensions
import App from './App.vue';
// import './index.css';
import 'buefy/dist/buefy.css';

// const main = () => {
//   API.test('index');
// };

// if (typeof document !== 'undefined') {
//   document.addEventListener('DOMContentLoaded', main);
// }

// use individually
// https://buefy.org/documentation/start
Vue.use(Buefy);
const vm = new Vue({ render: (h) => h(App) });
vm.$mount('#app');

export default {};
