import {readable, writable} from 'svelte/store';

export const battles = writable([]);

export const activeBattle = writable(undefined);

export const iteration = writable({
  ClientVersion: '',
  Ships: [],
});

// Wargaming API Information for a ship
export const shipInfo = writable({});

export const clientId = readable('2ecce5b4b0ffcffc5e7bc04131fb5c8e');

export const darkMode = writable(true);
