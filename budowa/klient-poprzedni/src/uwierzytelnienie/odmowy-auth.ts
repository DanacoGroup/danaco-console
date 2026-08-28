import { opisOdmowyBledu, type PowodOdmowy } from '../komponenty/odmowa';

/** Odmowa trójczęściowa bramki nazywa co się nie udało, dlaczego i czym operator to zmieni; rada domyślna odsyła do dziennika rdzenia, gdy kod odmowy nie ma osobnego zdania. */
export const RADA_OGOLNA = 'Powtórz czynność; jeśli odmowa wraca, sprawdź dziennik rdzenia.';

/** Rada na odpowiedź bez wpiętej rodziny komend bramki w rdzeniu, gdy montaż uwierzytelniania jej brakuje. */
const RADA_BEZ_UCHWYTU =
  'Rdzeń odpowiedział kopertą auth.unknown — rodzina auth.* nie jest wpięta do tego rdzenia. ' +
  'Uruchom rdzeń zbudowany z montażem bramki i odczytaj ekran ponownie.';

const RADY: Record<string, Record<string, string>> = {
  wejscie: {
    // Hasłu niezgodnemu rdzeń odmawia jednym kodem, a brakom kształtu żądania innym kodem.
    not_authenticated:
      'Hasło nie zgadza się z zapisem bramki — wpisz je ponownie. ' +
      'Jeżeli odmowa mówi o oczekiwaniu na potwierdzenie adresu, hasło jest dobre: ' +
      'przepisz drogę z listu odnośnikiem „Mam drogę potwierdzenia z listu".',
    validation_failed: 'Wpisz hasło — bramki bez hasła otworzyć się nie da.',
    not_found:
      'Bramki jeszcze nie ustawiono — ekran przechodzi na ustawienie hasła pierwszego uruchomienia.',
    internal_error:
      'Rdzeń pracuje bez sejfu poświadczeń albo sejf nie niesie sekretu — sprawdź dziennik rdzenia; ' +
      'bez sejfu bramka nie działa i żadne hasło nie wejdzie.',
  },
  // Wejście PIN-em ma własny obszar, bo te same kody znaczą tu co innego niż przy haśle.
  'wejscie-pin': {
    not_authenticated: 'PIN nie zgadza się z zapisem tej maszyny — wpisz go ponownie.',
    validation_failed:
      'Wpisz PIN — wejścia bez PIN-u rdzeń nie przyjmie. ' +
      'Jeśli odmowa mówi o urządzeniu, pamięć przeglądarki jest niedostępna i tożsamość maszyny ' +
      'nie ma się gdzie zapisać; wejdź hasłem.',
    not_found:
      'PIN nie jest założony na tej maszynie — wejdź hasłem, a PIN założysz w Ustawieniach, ' +
      'w sekcji „Uwierzytelnianie".',
    internal_error:
      'Rdzeń pracuje bez sejfu poświadczeń albo sejf nie niesie sekretu — sprawdź dziennik rdzenia; ' +
      'bez sejfu bramka nie działa i żaden PIN nie wejdzie.',
  },
  zalozenie: {
    validation_failed:
      'Wypełnij login, adres e-mail i hasło — rejestracji bez żadnego z nich rdzeń nie przyjmie.',
    conflict:
      'Konto właściciela już istnieje — rejestracja wykonuje się raz. Wpisz hasło w formularzu ' +
      'wejścia, a jeżeli go nie pamiętasz, odzyskaj konto odnośnikiem „Nie pamiętam hasła".',
    not_authenticated:
      'Konto czeka na potwierdzenie adresu — przepisz drogę potwierdzenia z listu. ' +
      'Formularz otworzy odnośnik „Mam drogę potwierdzenia z listu" pod polami.',
    internal_error:
      'Platforma nie ma konta nadawczego, więc listu z drogą potwierdzenia nie ma czym nadać. ' +
      'Ustaw serwer poczty wychodzącej: przy starcie rdzenia zmiennymi DANACO_NADAWCA_*, ' +
      'a po wejściu w oknie Konfiguracji, w kategorii „Konto nadawcze platformy". ' +
      'Ten sam kod niesie rdzeń pracujący bez sejfu poświadczeń — odmowa mówi wprost, o co chodzi.',
  },
  // Krok drugi rejestracji: droga z listu w zamian za token dostępu; każda odmowa dotyczy samej drogi.
  potwierdzenie: {
    validation_failed:
      'Przepisz drogę potwierdzenia z listu — pole nie może zostać puste. ' +
      'Droga jest jednorazowa i wygasa po godzinie.',
    not_found:
      'Tej drogi platforma nie zna — sprawdź, czy przepisałeś ją w całości i z właściwego listu. ' +
      'Droga zużyta albo starsza niż godzina już nie działa: poproś o nową, zaczynając odzyskanie ' +
      'konta odnośnikiem „Nie pamiętam hasła".',
    conflict:
      'Droga wygasła albo została już zużyta — poproś o nową odnośnikiem „Nie pamiętam hasła"; ' +
      'przyjdzie listem na ten sam adres.',
    not_authenticated:
      'Droga została wydana do innej czynności niż potwierdzenie adresu — użyj drogi z listu ' +
      'rejestracyjnego albo poproś o nową.',
    internal_error:
      'Rdzeń nie mógł domknąć potwierdzenia — sprawdź dziennik rdzenia; ' +
      'konto zostaje niepotwierdzone, a droga z listu jest nadal ważna do godziny od nadania.',
  },
  // Obie strony odzyskania konta: prośba o list i ustawienie nowego hasła drogą z listu.
  odzyskanie: {
    validation_failed:
      'Podaj adres e-mail uwierzytelniający konta, a w kroku drugim drogę z listu i nowe hasło — ' +
      'żadne z tych pól nie może zostać puste.',
    not_found:
      'Tej drogi platforma nie zna — przepisz ją z listu w całości. ' +
      'Droga zużyta albo starsza niż godzina już nie działa: poproś o nową.',
    conflict: 'Droga wygasła albo została już zużyta — poproś o nową i przepisz ją z nowego listu.',
    not_authenticated:
      'Droga została wydana do innej czynności niż odzyskanie konta — użyj drogi z listu ' +
      'o odzyskaniu albo poproś o nową.',
    internal_error:
      'Platforma nie ma konta nadawczego, więc listu z drogą odzyskania nie ma czym nadać. ' +
      'Ustaw serwer poczty wychodzącej: przy starcie rdzenia zmiennymi DANACO_NADAWCA_*, ' +
      'a po wejściu w oknie Konfiguracji, w kategorii „Konto nadawcze platformy".',
  },
  zmiana: {
    // Zmiana hasła odmawia inaczej niż logowanie; zdanie stoi, by wyrównanie kodów nie zostawiło rady.
    not_authenticated:
      'Dotychczasowe hasło nie zgadza się z zapisem bramki — wpisz je ponownie. ' +
      'Hasła zapomnianego ta droga nie odzyska, bo pyta o hasło dotychczasowe; ' +
      'od tego jest odzyskanie konta listem — odnośnik „Nie pamiętam hasła".',
    validation_failed:
      'Dotychczasowe hasło nie zgadza się z zapisem bramki albo któreś z pól zostało puste — ' +
      'wpisz oba hasła ponownie. ' +
      'Hasła zapomnianego ta droga nie odzyska, bo pyta o hasło dotychczasowe; ' +
      'od tego jest odzyskanie konta listem — odnośnik „Nie pamiętam hasła".',
    not_found:
      'Hasła jeszcze nie ustawiono na tym komputerze — wróć do wejścia i ustaw je przy pierwszym logowaniu.',
    internal_error:
      'Rdzeń pracuje bez sejfu poświadczeń — nowego hasła nie ma gdzie odłożyć; sprawdź dziennik rdzenia.',
  },
  przedluzenie: {
    conflict:
      'Poprzednia sesja wygasła albo została unieważniona — wejdź hasłem na nowo; ' +
      'praca w toku nie została przerwana, bo wygaśnięcie działa tylko przy wejściu.',
    not_found: 'Rdzeń nie zna zapisanej sesji — wejdź hasłem na nowo.',
    validation_failed: 'Zapis sesji był pusty — wejdź hasłem na nowo.',
    internal_error:
      'Rdzeń nie mógł odczytać zapisu sesji — sprawdź dziennik rdzenia; ' +
      'wejście hasłem działa niezależnie od tej odmowy.',
  },
  rozpoznanie: {
    internal_error:
      'Rdzeń nie może rozstrzygnąć stanu bramki — najpewniej pracuje bez sejfu poświadczeń; ' +
      'sprawdź dziennik rdzenia i odczytaj ekran ponownie.',
  },
};

/** Składa pełne zdanie odmowy: co się nie udało, dlaczego wprost z rdzenia, i czym operator to zmieni według obszaru czynności. */
export function opisOdmowyBramki(
  czynnosc: string,
  obszar: string,
  blad?: PowodOdmowy | null,
  bezUchwytu = false,
): string {
  const kod = (blad?.code ?? '').trim();
  const rada = bezUchwytu ? RADA_BEZ_UCHWYTU : (RADY[obszar]?.[kod] ?? RADA_OGOLNA);
  return `${opisOdmowyBledu(czynnosc, blad)}. ${rada}`;
}
