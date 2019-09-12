import axios from 'axios';

export function getShipInfo(id) {
  return axios.get(`https://api.worldofwarships.eu/wows/encyclopedia/ships/?application_id=2ecce5b4b0ffcffc5e7bc04131fb5c8e&ship_id=${id}`).then(res => res.data.data[`${id}`]);
}

export const shipClassImages = {
  'Destroyer': 'https://glossary-wows-global.gcdn.co/icons//vehicle/types/Destroyer/normal_578f4cb5e5bd1007d5df8ade07a8c4836e98436c2163927b1d14379ed318fec6.png',
  'AirCarrier': 'https://glossary-wows-global.gcdn.co/icons//vehicle/types/AirCarrier/normal_a103bcc3490cf6e21839ad78b2ee484bc2c8079f379526a71a8ef05354fd338c.png',
  'Battleship': 'https://glossary-wows-global.gcdn.co/icons//vehicle/types/Battleship/normal_57a69d35d27f32ad4c2ad55e76f684f474608e2f49162663c8574efc14be6f9f.png',
  'Cruiser': 'https://glossary-wows-global.gcdn.co/icons//vehicle/types/Cruiser/normal_50015a4b81edf8dec325a3c030cb526dc58b27a893ff2c6d25d4569a1ffe3775.png',
  'Submarine': 'https://glossary-wows-global.gcdn.co/icons//vehicle/types/Submarine/normal_658c33c080e9964038c12fa2e76bbfc264e007cd7d251e95e5f0159606072d71.png',
};

export const nations = {
  'commonwealth': 'Commonwealth',
  'europe': 'Europe',
  'italy': 'Italy',
  'usa': 'U.S.A.',
  'pan_asia': 'Pan-Asia',
  'france': 'France',
  'ussr': 'U.S.S.R.',
  'germany': 'Germany',
  'uk': 'U.K.',
  'japan': 'Japan',
  'pan_america': 'Pan-America',
};

export const tiers = ['I', 'II', 'III', 'IV', 'V', 'VI', 'VII', 'VIII', 'IX', 'X'];

