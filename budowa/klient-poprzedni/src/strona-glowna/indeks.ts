/** Plik jest punktem zbiorczym strony głównej: jedyne wejście dla warstw zewnętrznych, przez które powłoka aplikacji sięga po budowę strony głównej. */

export { utworzStroneGlowna, type StronaGlowna } from './okno-strona-glowna';
export { type KodSrodowiska, type PozycjaSrodowiska } from './pozycje-srodowisk';
export { type KodKomponentu, type PozycjaKomponentu } from './pozycje-komponentow';
export { pozycjaZKomponentu, pozycjeZKomponentow } from './pozycje-personalizowane';
export { type StrefaKomponentow } from './strefa-komponentow';
export { type StrefaModulow } from './strefa-modulow';
export { type PozycjaModuluStrony, pozycjeModulowBezNawigacji } from './pozycje-modulow';
export { wepnijModuly, type ZaleznosciWpieciaModulow } from './wpiecie-modulow';
export { type KodUstawienia, type PozycjaUstawienia } from './pozycje-ustawien';
export { type SluchaczWyboru } from './sygnal-wyboru';
export { wepnijSrodowiska } from './wpiecie-srodowisk';
export {
  wepnijKomponenty,
  type WpiecieKomponentow,
  type ZaleznosciWpieciaKomponentow,
} from './wpiecie-komponentow';
export { utworzMacierzModulow, type MacierzModulow } from './macierz-modulow';
export { zasilSesjeStronyGlownej, type ZaleznosciWpieciaSesji } from './wpiecie-sesji';
export { type MigawkaSesji, type WpisSesji } from './zrodlo-sesji';
