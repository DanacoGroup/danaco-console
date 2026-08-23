/**
 * Warstwa łączności — interfejs katalogu.
 *
 * Jedne drzwi do transportu, magistrali i obserwatorów zdarzeń. Kanał
 * kontraktu zakłada tu dziennik nierozpoznanych, a widoki sięgają po
 * obserwatora ogniska:
 *
 *   const dziennik = zalozDziennikNieznanych(kanal);
 *   const ognisko = utworzObserwatorOgniska(kanal);
 *   ognisko.naSesje(idSesji, (zmiana) => karta.ustawCzynna(zmiana.sessionId === idSesji));
 *
 * Obserwator osadza się na `ZrodloZdarzen`, nie na `Kanal` — dzięki temu
 * łączność nie zależy od warstwy protokołu i strzałka warstw pozostaje jedna.
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
