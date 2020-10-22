import API from '../api';

const main = () => {
  API.test('index');
};

if (typeof document !== 'undefined') {
  document.addEventListener('DOMContentLoaded', main);
}
