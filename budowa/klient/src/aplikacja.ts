// Punkt wejścia klienta. Znacznik okien pochodzi z biblioteki `design/zasoby/`.
import {
  adresGniazdaOdPowloki,
  adresGniazdaRdzenia,
  adresRdzeniaLokalnego,
  czyPochodzeniePowloki,
} from './polaczenie/adres-rdzenia.ts';
import { utworzTransport } from './polaczenie/gniazdo.ts';
import { utworzKanal } from './protokol/kanal.ts';
import { zadajPowitanie, zapomnijPowitanie } from './protokol/powitanie.ts';
import { sesjaKlienta } from './protokol/sesja.ts';
import { pilnujTokenu, tokenSesji } from './protokol/token-sesji.ts';
import { tozsamoscKlienta } from './protokol/tozsamosc-klienta.ts';
import { oglos } from './wiazanie/ogloszenie.ts';
import { pokazOknoPowloki, zwiazBelkeOkna } from './wiazanie/belka-okna.ts';
import { zwiazZmianeAdresu } from './wiazanie/konto-zmiana-adresu.ts';
import { zwiazWyglad } from './wiazanie/wyglad-okna.ts';
import { zwiazUrzadzenia } from './wiazanie/urzadzenia-konta.ts';
import { zwiazWejscie } from './wiazanie/wejscie.ts';
import { przejdzDoModulu } from './wiazanie/przejscie-do-modulu.ts';
import { zwiazZdarzenia } from './wiazanie/zdarzenia.ts';

/**
 * Adres gniazda rdzenia z dokumentu wczytanego po HTTP; poza tym pętla zwrotna
 * z portem domyślnym. Pochodzenie powłoki desktopowej do pętli zwrotnej nie
 * cofa — rdzeń stoi na serwerze wdrożenia, który podaje wskazanie powłoki.
 * Pustka znaczy powłokę bez wskazania: transport nie łączy, okno wejścia
 * pokazuje błąd.
 */
function adresRdzenia(): string | null {
  const pochodzenie = globalThis.location.origin;
  if (czyPochodzeniePowloki(pochodzenie)) return null;
  return adresGniazdaRdzenia(pochodzenie) ?? adresRdzeniaLokalnego();
}

/* Powłoka desktopowa pytana jest pierwsza: jej wskazanie niesie serwer
   wdrożenia, którego pochodzenie dokumentu w powłoce nie zdradza. Poza powłoką
   odpowiedzi nie ma i zostaje adres wywiedziony z pochodzenia. */
const transport = utworzTransport(adresGniazdaOdPowloki() ?? adresRdzenia() ?? '');
const kanal = utworzKanal(transport, sesjaKlienta());

/*
Kanał wystawiony na obiekcie globalnym: skrypty biblioteki są funkcjami
domkniętymi wczytywanymi przed tym modułem, nie modułami — importu z nich nie ma
i mieć nie może. Wystawienie jest jedynym stykiem między nimi a rdzeniem.
*/
declare global {
  // eslint-disable-next-line no-var
  var DanacoKanal: typeof kanal | undefined;
  // eslint-disable-next-line no-var
  var DanacoPrzejscieDoModulu: ((kodModulu: string, idSesji?: string) => void) | undefined;
}
globalThis.DanacoKanal = kanal;

/* Przejście z przedsionka do okna modułu wystawione globalnie: skrypty
   biblioteki i wiązania przedsionków wchodzą w moduł jedną drogą. */
globalThis.DanacoPrzejscieDoModulu = (kodModulu: string, idSesji = ''): void =>
  przejdzDoModulu(kanal, kodModulu, idSesji);

// Ekran startowy niesie `data-czekaj` i stoi do wywołania `gotowe()`.
function zglosGotowosc(): boolean {
  const pole = document.querySelector('[data-ekran-startowy]') as
    (HTMLElement & { ekranStartowy?: { gotowe?: () => void } }) | null;
  const gotowe = pole?.ekranStartowy?.gotowe;
  if (gotowe === undefined) return false;
  gotowe();
  return true;
}

/* Składnik ekranu startowego staje dopiero przy DOMContentLoaded, a pierwszy
   stan transportu potrafi paść wcześniej. Sygnał niedoręczony jest ponawiany —
   zgubiony zostawiałby ekran startowy czekający na `gotowe()` bez końca. */
let gotowoscZgloszona = false;
let ponowienie: ReturnType<typeof setTimeout> | undefined;
function raz(): void {
  if (gotowoscZgloszona) return;
  if (zglosGotowosc()) {
    gotowoscZgloszona = true;
    return;
  }
  if (ponowienie === undefined) {
    ponowienie = globalThis.setTimeout(() => {
      ponowienie = undefined;
      raz();
    }, 100);
  }
}

transport.naStan(raz);

/* Stan łączności jest widoczny dla Operatora. Warstwa projektowa nie niesie dla
   niego wzoru w oknie, więc mówi o nim komunikat biblioteki: zerwanie i powrót
   są zdarzeniami, o których praca musi wiedzieć, a milczące ponawianie wygląda
   jak program, który przestał odpowiadać. */
let bylPolaczony = false;
transport.naStan((stan) => {
  if (stan === 'polaczony') {
    if (bylPolaczony) oglos('Połączenie', 'Łączność z rdzeniem wróciła.');
    bylPolaczony = true;
    return;
  }
  if (!bylPolaczony) return;
  if (stan === 'rozlaczony') oglos('Połączenie', 'Łączność z rdzeniem zerwana.', 'ostrzezenie');
  if (stan === 'ponawianie') oglos('Połączenie', 'Wznawianie łączności z rdzeniem.', 'ostrzezenie');
});
globalThis.setTimeout(raz, 6000);

/* Token sesji bramki przejmowany jest z odpowiedzi rdzenia, bo to on rozstrzyga
   o ważności sesji; powitanie ponawiane niesie go z powrotem. */
pilnujTokenu(kanal);

/*
Powitanie idzie na każdym gnieździe od nowa: rdzeń wiąże sesję bramki
z połączeniem właśnie w powitaniu, więc komenda bez ponowienia wraca odmową
`not_authenticated`. Do powrotu odpowiedzi kolejka wychodząca stoi wstrzymana.
*/
transport.naStan((stan) => {
  /* Stan `laczenie` to pierwsza próba, przed którą żadnego połączenia nie było:
     powitanie zamówione wcześniej przez okno wejścia czeka na to samo gniazdo. */
  if (stan === 'rozlaczony' || stan === 'ponawianie') {
    zapomnijPowitanie();
    return;
  }
  if (stan === 'polaczony') void przywitaj();
});

let powitanoJuz = false;
async function przywitaj(): Promise<void> {
  const wynik = await zadajPowitanie(kanal, tozsamoscKlienta(), tokenSesji() || undefined);
  transport.zwolnijWstrzymanie();
  /* Nieudane powitanie pierwsze ma swój obraz w oknie wejścia — wariant błędu
     stoi tam zamiast okna logowania. Dalsze padają przy oknie już zamkniętym,
     więc bez komunikatu Operator patrzyłby w okno, które przestało odpowiadać. */
  if (!wynik.udany && powitanoJuz) {
    oglos('Połączenie', wynik.blad?.message ?? 'Rdzeń odrzucił powitanie.', 'blad');
  }
  powitanoJuz = true;
}

/* Zdarzenia rdzenia idą jednym rozdzielaczem: wiązania okien zgłaszają do niego
   uchwyty, a warstwa wspólna ma odbiorcę od chwili postawienia aplikacji. */
zwiazZdarzenia(kanal);

/* Wiązanie znacznika Właściciela z komendami rdzenia: nasłuchy na jego
   przyciskach i polach. Nie stawia żadnego elementu. */
zwiazWejscie(kanal);

/* Zmiana adresu konta: zamówienie i potwierdzenie z okna, wycofanie drogą
   z listu — ta ostatnia wykonuje się przy wejściu, bez kliknięcia. */
zwiazZmianeAdresu(kanal);

/* Motyw i gęstość należą do konta, nie do przeglądarki: wybór idzie do rdzenia
   i wraca przy każdym otwarciu okna. */
zwiazWyglad(kanal);

/* Wykaz urządzeń bierze się z rdzenia: prototyp niesie przykłady, a instalacja
   ma tyle urządzeń, ile faktycznie się na niej uwierzytelniło. */
zwiazUrzadzenia(kanal);

/* Belka okna zastępuje ramę systemową, więc sterowanie oknem idzie z niej. */
zwiazBelkeOkna();

/* Okno powłoki wstaje ukryte i pokazuje się dopiero wtedy, gdy ekran startowy
   stoi zmontowany i narysowany — inaczej Operator patrzy przez kilka sekund
   w puste pole, zanim biblioteka zdąży złożyć okno. Zapora czasu jest po to,
   żeby usterka montażu nie zostawiła okna niewidocznym na zawsze. */
const ZAPORA_POKAZANIA_MS = 4000;
const ZAPORA_KLATKI_MS = 500;
function ekranStartowyStoi(): boolean {
  const pole = document.querySelector('[data-ekran-startowy]') as
    (HTMLElement & { ekranStartowy?: unknown }) | null;
  return pole?.ekranStartowy !== undefined;
}
/* Okno nigdy niepokazane nie jest składane, a silnik widoku nie wywołuje wtedy
   `requestAnimationFrame` — samo czekanie na klatkę zamyka więc okno w zapętleniu:
   klatka czeka na pokazanie, pokazanie na klatkę. Klatka zostaje drogą pierwszą,
   a zegar drogą zapasową; pierwsza z nich pokazuje okno, druga nie robi nic. */
let pokazano = false;
function pokazRaz(): void {
  if (pokazano) return;
  pokazano = true;
  pokazOknoPowloki();
}
function pokazPoZlozeniu(odKiedy: number): void {
  if (ekranStartowyStoi() || Date.now() - odKiedy >= ZAPORA_POKAZANIA_MS) {
    requestAnimationFrame(() => requestAnimationFrame(pokazRaz));
    globalThis.setTimeout(pokazRaz, ZAPORA_KLATKI_MS);
    return;
  }
  globalThis.setTimeout(() => pokazPoZlozeniu(odKiedy), 50);
}
pokazPoZlozeniu(Date.now());

transport.polacz();
