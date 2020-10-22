import API from '../api';

const main = () => {
  API.test('get');
};

if (typeof document !== 'undefined') {
  document.addEventListener('DOMContentLoaded', main);
}
