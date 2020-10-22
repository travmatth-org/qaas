import API from '../api';

const main = () => {
  API.test('lost');
};

if (typeof document !== 'undefined') {
  document.addEventListener('DOMContentLoaded', main);
}
