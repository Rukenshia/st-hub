import {writable} from 'svelte/store';

export const battles = writable([]);

export const iteration = writable({
  clientVersion: '',
  iterationName: '',
  ships: [],
});
