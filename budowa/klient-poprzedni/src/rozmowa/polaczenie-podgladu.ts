import { adresRdzenia } from '../polaczenie/adres-rdzenia';
import { utworzTransport } from '../polaczenie/gniazdo';
import { opisPoczatkowy } from '../okno-komunikacji/opis-okna';
import { zamowienieOkna } from '../okno-komunikacji/zamowienie-okna';
import { utworzKanal } from '../protokol/kanal';
import { utworzSesje } from '../protokol/sesja';
import { tozsamoscKlienta } from '../protokol/tozsamosc-klienta';
import { utworzUzgodnienie } from '../protokol/uzgodnienie';
import { zamontujRozmowe } from './indeks';
import { zapewnijKanalGlowny } from './zapewnienie-kanalu';

/**
 * Miejsca w interfejsie podglądu, w których uruchomiony podgląd wypisuje przebieg łączenia
 * z rdzeniem.
 */
export interface CzesciPodgladu {
  /** Kontener widoku rozmowy. */
  scena: HTMLElement;
  /** Wiersz stanu nad sceną. */
  stan: HTMLElement;
}

/**
 * Uruchamia podgląd warstwy rozmowy na żywym rdzeniu tą samą drogą, którą przechodzi
 * powłoka produktu.
 */
export function uruchomPodglad(czesci: CzesciPodgladu): void {
  const opis = opisPoczatkowy();
  const transport = utworzTransport(adresRdzenia());
  const kanal = utworzKanal(transport, utworzSesje());
  let rozpoczete = false;

  function zglos(tekst: string): void {
    czesci.stan.textContent = tekst;
  }

  transport.naStan((stanPolaczenia) => {
    if (stanPolaczenia !== 'polaczony') {
      zglos(`Rdzeń ${adresRdzenia()} — ${stanPolaczenia}`);
      return;
    }
    if (rozpoczete) return;
    rozpoczete = true;
    zglos('Połączono. Zapewniam kanał główny w rejestrze…');
    zapewnijKanalGlowny(kanal, (wynik) => {
      if (wynik.idKanalu.length === 0) {
        zglos(`Rejestr kanałów niedostępny — ${wynik.przeszkoda}`);
        return;
      }
      zglos(`Kanał główny ${wynik.pochodzenie}: ${wynik.idKanalu}. Otwieram okno…`);
      otworzOkno(wynik.idKanalu);
    });
  });

  function otworzOkno(idKanalu: string): void {
    const uzgodnienie = utworzUzgodnienie(
      kanal,
      { ...zamowienieOkna(opis), modelChannelId: idKanalu },
      tozsamoscKlienta(),
    );
    uzgodnienie.naPostep((postep) => {
      if (postep.udany) return;
      zglos(`Uzgodnienie — etap ${postep.etap} nieudany: ${postep.blad?.message ?? 'brak przyczyny'}`);
    });
    uzgodnienie.naOtwarcieOkna((okno) => {
      zglos(`Okno ${okno.id} · kanał ${idKanalu} · rola ${okno.windowRole}`);
      const zamontowana = zamontujRozmowe(czesci.scena, kanal, okno.id, {
        persona: opis.kanalModelu,
        rolaOkna: okno.windowRole,
      });
      const pytanie = pytanieZAdresu();
      if (pytanie.length > 0) zamontowana.rozmowa.wyslij(pytanie);
    });
    uzgodnienie.rozpocznij();
  }

  transport.polacz();
}

/**
 * Wypowiedź podana w adresie podglądu (`?pytanie=…`).
 *
 * Pozwala pokazać pełną turę — prowenancję, strumień, narzędzia
 * i podsumowanie — bez ręcznego pisania. Brak parametru zostawia podgląd
 * czekający na wypowiedź Operatora.
 */
function pytanieZAdresu(): string {
  const lokalizacja = globalThis.location;
  if (!lokalizacja) return '';
  return new URLSearchParams(lokalizacja.search).get('pytanie') ?? '';
}
