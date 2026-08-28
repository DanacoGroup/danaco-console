/**
 * Warstwa łączności jest interfejsem katalogu: jedne drzwi do transportu, magistrali i obserwatorów zdarzeń, gdzie obserwator osadza się na źródle zdarzeń, nie na kanale, więc łączność nie zależy od warstwy protokołu.
 */
export { adresRdzenia, przyjmijAdresZPowloki } from './adres-rdzenia';

export { utworzTransport, type Transport } from './gniazdo';

export { utworzKolejkeWychodzaca, type KolejkaWychodzaca } from './kolejka-wychodzaca';

export {
  utworzMagistrale,
  type Magistrala,
  type Odsubskrybuj,
  type Sluchacz,
} from './magistrala-zdarzen';

export { wykladniczePonawianie, type PolitykaPonawiania } from './ponawianie';

export type { StanPolaczenia } from './stan-polaczenia';

export type { ZrodloZdarzen } from './zrodlo-zdarzen';

export { utworzObserwatorOgniska, type ObserwatorOgniska } from './obserwator-ogniska';

export {
  zalozDziennikNieznanych,
  type DziennikNieznanych,
  type WpisNieznanego,
} from './dziennik-nieznanych';
