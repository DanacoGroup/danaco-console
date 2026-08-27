# Uzasadnienia pakietów

| | |
|---|---|
| **Produkt** | Danaco Console |
| **Przeznaczenie** | Uzasadnienia projektowe i zastrzeżenia przeniesione z nagłówków kodu, których pełna treść przekraczała dopuszczalną długość komentarza |

Dokument gromadzi rozdziały nazwane ścieżką pliku źródłowego. Każdy rozdział wylicza identyfikatory (typ, pole, stała, funkcja) wraz z uzasadnieniem, które nie mieściło się w nagłówku komentarza tego pliku.

---

### budowa/shared/contract.go

Uzasadnienia i zastrzeżenia projektowe przeniesione z komentarzy kontraktu komunikacji WebSocket, którym nagłówek kodu nie miał miejsca.

**MessageType** — Jest aliasem napisu, wiec stale kontraktu wchodza wprost tam, gdzie warstwa protokolu operuje na string.

**ExecutionEnv** — Nazwa konkretnego hosta nie jest wartoscia wyliczenia — idzie ustawieniem na poziomie zasiegu window

**PermissionMode** — Wartosci odpowiadaja doslownie przelacznikowi --permission-mode kanalu glownego; zweryfikowane z claude --help 2026-08-11.

**ChunkKind** — Ten sam slownik opisuje rodzaj tresci wiadomosci zapisanej w bazie

**ChunkKind.ChunkKindFinal** — ZASTEPUJE tresc zlozona z dotychczasowych fragmentow tekstowych, nie doklada sie do niej. Niesie ja koperta domykajaca (done=true); odbiorca podmienia tekst wpisu na tresc tego fragmentu. Wartosc PRZELOTOWA — tresc ostateczna utrwala sie jako `content` wiadomosci o rodzaju `tekst`, wiec sam ten rodzaj do kolumny nigdy nie trafia

**ConfigScope.ConfigScopeApplication** — Poziom NAJSZERSZY: przegrywa z kazdym wezszym zapisem, a zapisu na nim nie ma czym zawezic, bo bytu nie ma (scopeId pusty). Tedy przelacznik ustalenie budowy dochodzi do warstwy nasluchu bez restartu rdzenia

**NavigationKind** — Trzy srodowiska wystawiaja liste modulow, MultitaskingAI wystawia panel orkiestracji i nie udostepnia modulow

**ConfigAxis** — Os jest prostopadla do poziomu zasiegu: poziom mowi jak wasko (okno -> ... -> globalny), os mowi dla czego (platforma, model, konto). Klucz rozstrzygania jest zlozony: klucz + poziom + byt poziomu + os + byt osi. Brak osi znaczy platform

**ConfigAxis.ConfigAxisPlatform** — Wartosc domyslna osi — os pominieta znaczy platform

**SettingValueType** — Klient buduje z niego kontrolke formularza i nie zna zadnego klucza z osobna

**AccessPointKind** — Punkt dostepu nie jest srodowiskiem: srodowisko pozostaje profilem widocznosci modulow w bocznej nawigacji

**IdentityLayer** — Wartosci odpowiadaja stalym silnika nakladki

**IdentityMode** — Tozsamosc jest ZAMIENIANA, nie dolaczana — wartoscia domyslna klucza tozsamosc.tryb_domyslny jest ZASTAP. Wartosci sa sterowaniem istniejacego silnika nakladki, nie jego powtorzeniem

**IdentityMode.IdentityModeZASTAP** — Tryb zastapienia nie podaje --append-system-prompt

**LoopStopReason** — Zatrzymanie nigdy nie jest ciche; wartosci odpowiadaja stalym session.PowodZatrzymania

**LoopStopReason.LoopStopReasonCompleted** — Jedyna wartosc mowiaca skonczone, a nie przerwane — po niej bieg podejmuje sie sam przy kolejnym koncu tury wykonawcy, wiec ukonczenie nie jest brama akceptacji; odpowiada session.ZatrzymanieUkonczenie

**SessionConfigArea** — Obszar jest jednostka zapisu: komenda config.session.set bierze wykaz obszarow, nie wykaz pol. Warstwy tozsamosci, profilu i ekspertyzy nie sa obszarami — sa polami obszaru systemPrompt

**ProviderTransport** — Wartosc wyznacza adapter, ktory tlumaczy konfiguracje sesji na powierzchnie dostawcy — sam model konfiguracji pozostaje ten sam

**AdapterSurface** — Interfejs nie jest budowany wokol nazw powierzchni — sluza wylacznie diagnostyce i deklaracji zdolnosci

**SessionConfigSource** — Konfiguracja sesji SPINA istniejace rodziny, nie zastepuje ich — pole mowi, ktory rejestr byl zrodlem prawdy

**WorkingDirectoryDegradation** — Rozejscie ma byc widoczne w oknie, nie milczace

**DesignAssetKind.DesignAssetKindDocument** — Wartosc dopisana 15.08.2026 — narzedzia dokumentowe wytwarzaly pliki, ktore trzeba bylo opisywac jako wektor, bo enum nie znal dokumentu

**DesignAssetKind.DesignAssetKindAudio** — Wartosc dopisana 15.08.2026 — narzedzia mediow wytwarzaly dzwiek, ktory trzeba bylo opisywac jako obraz

**DesignAssetKind.DesignAssetKindVideo** — Wartosc dopisana 15.08.2026 z tego samego powodu co audio

**AgentPermissionGroup.AgentPermissionGroupModules** — Zakres szczegolowy wpisu niesie kod modulu z rodziny `module.*`, wiec wykaz modulow nie jest tu powielany

**MemoryLevel** — WYLACZENIE PAMIECI NIE JEST PIATA WARTOSCIA TEGO WYLICZENIA — jest PUSTYM ZBIOREM poziomow w polu `Agent.memoryLevels`. Piata wartosc `disabled` obok czterech poziomow pozwalalaby zapisac stan sprzeczny (`[session, disabled]`), ktorego rdzen musialby rozstrzygac zgadywaniem; zbior pusty jest jedynym ksztaltem, w ktorym "wylaczona" nie da sie postawic obok "wlaczona". Wartosc wyjsciowa definicji eksperta jest nadpisywana przy przypisaniu do projektu i przy przypisaniu do roli — pierwszenstwo ma zasieg najbardziej szczegolowy

**AgentVisibility** — Ekspert `project` NIE MA wlasnego wskazania projektu: widoczny jest w tych projektach, do ktorych zostal przypisany komenda `workspace.agent.assign` — drugie miejsce zapisu przynaleznosci byloby druga prawda o tym samym

**ExtensionOrigin** — RDZEN NIE ROZROZNIA ZRODLA: ladowanie, kontrakt integracji i uzycie operacyjne sa dla obu zrodel te same. Jedyny wyjatek zapisany przez ustalenie budowy to STAN WYJSCIOWY PRZY REJESTRACJI — `danaco` staje wlaczone, `personal` wylaczone. Poza tym jednym miejscem zrodlo jest faktem do pokazania Operatorowi, nie przeslanka rozgalezienia zachowania

**AuthMethodKind** — Trzy, wzorem Danaco HUB: haslo Operatora, PIN i Windows Hello. Bez OAuth, bez federacji, bez rotacji kluczy, bez MFA. Konta uzytkownika NIE MA — jedyny Operator jest bezimienny

**AuthMethodKind.AuthMethodKindHello** — NIEDOSTEPNA POD BIEZACYM POCHODZENIEM: WebAuthn wywodzi rp_id z pochodzenia dokumentu, a powloka natywna podaje interfejs z http://127.0.0.1:<port>, czyli z adresu IP, ktory domena nie jest. Metoda staje sie dostepna dopiero, gdy interfejs jest podawany z pochodzenia z domena po https

**AdvisorSelection** — Dwie wartosci, bo drogi sa dwie: wskazanie Operatora albo sufit sily. Trzeciej — eskalacji z inicjatywy modelu — nie ma i nie bedzie: rozstrzygniecie ustalenie budowy z 14.08.2026 mowi, ze bez wskazania Operatora model kontraktuje agenta o parametrach takich samych albo nizszych

**AdvisorSelection.AdvisorSelectionOperatorIndication** — Wskazanie znosi sufit sily, bo pochodzi od czlowieka, nie od modelu

**AdvisorSelection.AdvisorSelectionStrengthCeiling** — Tak konczy sie kazda konsultacja, o ktora poprosil sam model

**ActorKind** — Bez tego rozroznienia Operator widzi skutek, ale nie widzi REKI — a asystent główny prace za niego staje sie nieodrozniahy od niego samego

**SlashEntryKind** — To NIE SA dwa mechanizmy: jedno pole, jeden wykaz, jedno filtrowanie — rozne sa skutki wyboru. Klasyfikacja musi istniec w danych, bo pozycja rodzaju `tool` POSZERZA zestaw narzedzi sesji, a pozycja rodzaju `action` wykonuje czynnosc aplikacji i zestawu nie dotyka

**SessionToolSource** — Komenda po ukosniku jest JEDYNA droga poszerzenia zestawu w trakcie pracy (rozdz. 2 „Wylacznosc"), a to pole mowi, kto te droge przeszedl — Operator z klawiatury czy asystent dzialajacy za niego. Nie jest to druga droga dolozenia, tylko zapis sprawcy

**AgentPermissionAction** — Dzisiejszy `AgentPermission` niesie grupe, zakres i wartosc logiczna, wiec zawezenie „konektor tylko do odczytu" nie ma jak zostac zapisane — wartosc logiczna rozstrzyga wylacznie „wolno albo nie wolno"

**AgentModuleScope** — Wartosci odpowiadaja przykladom wymienionym w opracowaniu rozdz. 9.2): zapis plikow w module Library, wykonywanie polecen w module Terminal

**TelemetryFormat** — Rozlaczny z ExportFormat, ktory sluzy dokumentom Operatora (pdf, docx, markdown) i formatow maszynowych nie niesie

**UsageDimension** — Jeden wymiar na zadanie: zestawienie po dwoch wymiarach naraz jest dwoma zadaniami, nie jednym

**AlertMetric** — Kazda wartosc ma zrodlo w danych, ktore platforma juz przechowuje — reguly nad miara bez zrodla nie da sie ewaluowac

**AlertChannel** — Poza droga w aplikacji kazda wymaga skonfigurowanej integracji; brak integracji nie wstrzymuje wyzwolenia, tylko dostarczenie

**LibraryFileStatus** — Archiwizacja jest odwracalna: zasob zarchiwizowany znika z wykazu domyslnego, nie z repozytorium

**LibraryRuleKind** — Obie skladaja plik z kolekcja, ale z innej strony: regula kolekcji liczy zawartosc z warunku, regula naplywu wciaga zasob z zewnatrz

**LibraryPreservationKind** — Kazda jest odrebnym standardem branzowym, nie wariantem jednego zapisu

**LibraryDuplicateKind** — Bez niej wynik dokladny i wynik przyblizony wygladaja tak samo

**LibraryFieldKind** — Rozstrzyga postac kontrolki i sprawdzian wartosci, nie blokade zapisu

**LibraryRetentionAction** — Twarde usuniecie nie zachodzi samoczynnie

**WorkspaceTaskStatus** — Kolumny tablicy sa nastawa Operatora (rozdz. 9 opracowania) i odwzorowuja sie na te stany; kolumna wlasna nie zaklada stanu nowego

**WorkspaceDependencyKind** — Dwie wartosci, bo tylko dwie ma pokrycie w opracowaniu: nastepstwo poprzednik-nastepnik oraz blokada

**WorkspaceEntityKind** — Jedno wyliczenie dla wyszukiwania, komentarzy, wezlow grafu i osi czasu — te same byty pod jedna nazwa, bez czterech list synonimow

**StudioActionKind.StudioActionKindFieldChange** — Osobno od objectChange, bo pole odlozone jako zmiana obiektu nie da sie cofnac samo — Operator cofalby wtedy cala zmiane obiektu razem z nia

**AodEventClass** — Klasa jest jednoczesnie zakresem wyciszenia (rozdz. 3.5, wiersz trzeci): wyciszona klasa nie tworzy sugestii do chwili zniesienia, a pozostale dzialaja bez zmian

**AodMuteKind** — Tryb cichy NIE JEST czwarta wartoscia tego wyliczenia: jest trybem obecnosci (rozdz. 9.2) i jedzie ustawieniem konfiguracji, nie wpisem wyciszenia. Wyjatek wagi krytycznej takze nie jest wyciszeniem — jest regula przebijajaca kazde z trzech

**AodMuteKind.AodMuteKindTimed** — Chwile konca niesie pole endsAt

**AodMuteScope** — Wyciszenie czasowe obowiazuje cala platforme i zakresu nie niesie (pole scope puste); wyciszenie kontekstowe wskazuje modul albo karte sesji, wyciszenie klasy zdarzen — klase

**MessageType.CommandSpeechAvailabilityGet** — Brak silnika to ODPOWIEDZ, nie awaria — klient nie pokazuje mikrofonu, ktory nic nie nagra

**MessageType.CommandSpeechTranscribe** — Dzwiek NIE OPUSZCZA maszyny, na ktorej stoi silnik

**MessageType.CommandSessionFocus** — Asystent przestawia ognisko OPERATOROWI polem targetClientId — inaczej jego posuniecia nie byly widoczne na ekranie

**MessageType.CommandDesignAssetGenerate** — Tresc trafia do magazynu rdzenia pod suma kontrolna, wiec zasob nie zalezy od zadnego pliku zewnetrznego. Brak kanalu obrazowego albo brak poswiadczenia to ODMOWA NAZYWAJACA BRAK — rdzen nigdy nie zaklada zasobu bez bajtow obrazu

**MessageType.CommandAgentLayerSet** — Warstwa DOPISUJE sie do promptu systemowego, nie zastepuje go — trybu podania sie tu nie wybiera

**MessageType.CommandMemoryDetach** — Odpiecie NIE JEST wylaczeniem: zweza zasieg samego wpisu, a nie wstrzymuje go w cudzym projekcie ani module — do tego sluzy memory.disable.set

**MessageType.CommandComponentAssign** — UWAGA: znaczenie przypisania nie zostalo rozstrzygniete przez ustalenie budowy

**MessageType.CommandExtensionInstall** — UWAGA: znaczenie instalacji — pobranie paczki, zarejestrowanie adresu czy zapis punktu dostepu — nie zostalo rozstrzygniete przez ustalenie budowy

**MessageType.CommandExtensionUninstall** — UWAGA: znaczenie odinstalowania jest zwiazane z nierozstrzygnietym znaczeniem instalacji

**MessageType.CommandAdvisorConsult** — Rada NIE JEST WIAZACA i nie wykonuje pracy za pytajacego — odpowiedzialnosc za wynik zostaje przy agencie, ktory pyta. Kim jest doradca, rozstrzyga RDZEN z danych okna, nie zadanie: pole requestedAdvisor jest PROSBA i podlega sufitowi sily. Prosba o model silniejszy od kanalu pytajacego bez wczesniejszego wskazania Operatora jest ODMAWIANA, a nie zamieniana po cichu na slabszego doradce — model z wlasnej inicjatywy po model silniejszy nie siega (rozstrzygniecie ustalenie budowy z 14.08.2026). Kazda odbyta konsultacja rozglasza sie zdarzeniem advisor.consulted

**MessageType.CommandConfigExplainGet** — Transparentnosc, nie bramka

**MessageType.CommandConfigWindowOpen** — UWAGA: wykaz okien obiecuje te komende, nie mowiac, co rdzen ma przy niej robic poza podaniem zakresu

**MessageType.CommandTerminalOutputRead** — Komenda `terminal.command.exec` konczy sie w chwili STARTU procesu (kompilacja trwa dluzej niz kazde sensowne oczekiwanie na odpowiedz, a rozlaczenie klienta nie ma prawa jej przerwac), wiec wyniku niesc nie moze i nigdy nie bedzie mogla. Tedy model czyta, co polecenie wypisalo. Zrodlem jest ten sam dziennik zbiorczego wyjscia, na ktorym stoi `terminal.output.stream` — drugiej pompy nie ma. Historia zyje jeden bieg rdzenia: po restarcie wyjscie jest puste i odpowiedz mowi to wprost

**MessageType.CommandModelChannelSet** — Kanaly zaklada i zmienia rodzina channel.*; ta komenda wybiera jeden z zalozonych

**MessageType.CommandAodSuggestion** — UWAGA: wykaz okien obiecuje te pozycje slowem podpowiedz, nie rozstrzygajac, czy jest to odczyt Operatora, czy zdarzenie wypychane przez rdzen

**MessageType.CommandAuthRegister** — Konto powstaje w stanie NIEPOTWIERDZONYM i pozostaje w nim do chwili potwierdzenia adresu komenda auth.verify — dopiero potwierdzenie wydaje urzadzeniu token dostepu. Sesji ta komenda NIE zaklada. Wykonalna tylko raz; potem odmawia trwale. Powtorzenie hasla jest sprawa formularza klienta, nie kontraktu

**MessageType.CommandAuthVerify** — Przenosi konto ze stanu niepotwierdzonego do potwierdzonego, po czym wydaje urzadzeniu token dostepu — to jest moment, w ktorym Operator wchodzi do platformy po raz pierwszy. Droga potwierdzenia przychodzi listem na adres podany przy rejestracji i wygasa; wygasla droga odmawia i pozwala poprosic o nowa

**MessageType.CommandAuthRecover** — Serwer wysyla na wskazany adres droge potwierdzenia tozsamosci; nowe haslo ustawia sie komenda auth.reset. Odpowiedz nie zdradza, czy adres pasuje do konta — inaczej komenda mowilaby obcemu, jaki adres ma Operator

**MessageType.CommandAuthReset** — Zastepuje dotychczasowy skrot hasla przechowywany poza baza danych i uniewaznia tokeny dostepu wydane przed odzyskaniem — kazde powiazane urzadzenie loguje sie ponownie, a zmiana idzie w swiat zdarzeniem device.changed. Nowej encji Konto nie tworzy: zmienia sie wylacznie material uwierzytelniajacy

**MessageType.CommandAuthLogin** — Metodami dzialajacymi dzis sa HASLO i PIN; hello czeka na pochodzenie z domena po https. Login niesie sie przy metodzie password, bo konto ma nazwe ustalenie budowy nadana przy rejestracji; metody wlasciwe urzadzeniu tozsamosci nie potrzebuja, bo wskazuje ja material na urzadzeniu. Odmowa wraca kodem bledu not_authenticated. Proba nieudana NIE odmawia nastepnej — nakłada na nia ZWLOKE rosnaca wykladniczo (250 ms, 500, 1000, 2000, 4000, dalej rowno 5000 ms), zerowana przy pierwszym udanym wejsciu. Operator, ktory pomylil haslo trzy razy, wchodzi za czwartym; tylko czeka. Zadnego progu prob i zadnej odmowy 'za duzo prob' tu nie ma i nie bedzie — to byloby bramkowanie (rozstrzygniecie ustalenie budowy z 14.08.2026)

**MessageType.CommandAuthMethodAdd** — Czynnosc USTAWIEN, wykonalna po zalogowaniu; hasla ta komenda nie zaklada, bo kotwica powstaje przy auth.register. Klucz Hello nie opuszcza TPM urzadzenia, wiec na kazdej maszynie zaklada sie go osobno — i tylko wtedy, gdy interfejs jest podawany z pochodzenia z domena po https. Dzis zalozyc mozna PIN

**MessageType.CommandAuthMethodRemove** — Ani ostatniej metody, ani hasla zdjac sie nie da — haslo jest kotwica bramki

**MessageType.CommandAuthTokenRefresh** — W Danaco HUB dzieje sie to samo przy kazdym zadaniu HTTP; tu jest komenda, bo transportem jest jedno gniazdo WebSocket

**MessageType.CommandDeviceList** — Konto jest jedno, urzadzen dowolnie wiele: komputery, telefony, tablety. Kazde niesie wlasny token dostepu, wiec wykaz jest miejscem, w ktorym Operator widzi, co ma dostep do platformy

**MessageType.CommandDeviceRevoke** — Zmiana idzie do wszystkich polaczonych urzadzen zdarzeniem device.changed, wiec Operator widzi skutek natychmiast na pozostalych ekranach. Tej samej drogi uzywa odzyskanie konta, ktore uniewaznia tokeny wydane przed zmiana hasla

**MessageType.CommandAppsDeploymentList** — Bez tej komendy historia wdrozen znikala po odswiezeniu okna, bo klient odbudowywal ja wylacznie ze zdarzen biezacej sesji

**MessageType.CommandAppsArchitectureGet** — Odpowiednik odczytu dla apps.architecture.define — bez niego okno po odswiezeniu nie wie, co Operator wczesniej zdefiniowal

**MessageType.CommandAppsWorkspaceList** — Odpowiednik odczytu dla apps.workspace.update

**MessageType.CommandDesignAssetUpload** — Tresc wchodzi do magazynu rdzenia pod suma kontrolna — zasob przestaje zalezec od pliku, ktory Operator moze nadpisac albo skasowac. Uzyj, gdy zasob JUZ ISTNIEJE: Operator ma plik albo wskazuje odsylacz. Zasob majacy dopiero powstac zamawia sie komenda design.asset.generate

**MessageType.CommandDesignAssetRemove** — Rodzaj zmiany deleted zdarzenia design.asset.changed nie mial dotad zadnego nadawcy

**MessageType.CommandImageTransform** — Wynikiem jest NOWY zasob w magazynie — zrodlo zostaje nietkniete, wiec model moze probowac bez ryzyka

**MessageType.CommandImageAdjust** — To jest RETUSZ, o ktory model prosi w rozmowie

**MessageType.CommandImageConvert** — Sluzy przygotowaniu zasobu do wydania — model uzywa jej przed osadzeniem grafiki na stronie

**MessageType.CommandMediaTranscode** — Wynikiem jest nowy zasob w magazynie

**MessageType.CommandDocumentConvert** — Model uzywa jej, gdy Operator prosi o dokument, a nie o tekst w oknie

**MessageType.CommandDocumentTextExtract** — Model uzywa jej, zeby PRZECZYTAC to, co dostal jako plik

**MessageType.CommandArchivePack** — Sluzy wydaniu pracy Operatorowi jednym plikiem

**MessageType.CommandMailAccountList** — Bez skonfigurowanej skrzynki rodzina mail.* odmawia, nazywajac brak — rdzen poczty nie zmysla

**MessageType.CommandMailMessageList** — Uzyj do odnalezienia sprawy, o ktorej mowi Operator —na przykład listu od konkretnej osoby z dzisiaj

**MessageType.CommandMailMessageGet** — Uzyj, gdy masz przeanalizowac list, a nie tylko go odnalezc. Zalaczniki trafiaja do magazynu rdzenia i wracaja jako zasoby, wiec model moze je dalej przetworzyc narzedziami obrazu i dokumentow

**MessageType.CommandMailDraftSave** — Szkic jest krokiem POSREDNIM: Operator widzi go w swojej poczcie, zanim cokolwiek wyjdzie w swiat

**MessageType.CommandMailSend** — TO JEST JEDYNA KOMENDA RODZINY, KTOREJ SKUTEK WYCHODZI POZA MASZYNE OPERATORA i nie da sie go cofnac

**MessageType.CommandMailAccountAdd** — Rdzen nie zaklada konta pocztowego i nie stawia serwera — bierze skrzynke, ktora juz istnieje: na urzadzeniu albo w chmurze. Poswiadczenie idzie do sejfu, nie do bazy

**MessageType.CommandMailAccountDiscover** — Niczego nie podpina i nie siega po haslo — oddaje to, co znalazl, zeby Operator nie przepisywal nastaw recznie

**MessageType.CommandMailAccountRemove** — Skrzynki u dostawcy nie tyka

**MessageType.CommandImageUpscale** — Zastepuje wyspecjalizowane modele powiekszajace; bez zainstalowanego silnika odmawia, nazywajac brak — nigdy nie oddaje zwyklego rozciagniecia jako powiekszenia

**MessageType.CommandImageBackgroundRemove** — Zastepuje wyspecjalizowane narzedzia wycinania; bez silnika odmawia, nazywajac brak

**MessageType.CommandKnowledgeIndex** — Wyszukiwanie po slowach juz dziala (library.file.search); to jest droga do wyszukiwania po SENSIE

**MessageType.CommandKnowledgeImageSearch** — Rodzina knowledge.* prowadzi dotad wylacznie tekst; ta komenda jest jej osia obrazu i przeglada obrazy biblioteki Operatora

**MessageType.CommandSessionToolAttach** — Definicji eksperta NIE RUSZA (rozstrzygniecie ustalenie budowy, rozdz. 5): dolozenie zyje w stanie sesji, przezywa rozlaczenie klienta i konczy sie wraz z sesja albo z usunieciem rozmowy

**MessageType.CommandSessionToolDetach** — Zestaw wraca do podstawy z definicji eksperta; sama definicja nie zmienia sie ani o joto, bo nigdy nie byla zmieniana

**MessageType.CommandSessionToolList** — Zestaw narzedzi tury to definicja eksperta PLUS te dolozenia — bez tego odczytu druga polowa zestawu bylaby niewidoczna

**MessageType.CommandToolsCatalogList** — Wykaz liczy setki pozycji, wiec zadanie niesie zawezenie tekstem, rodzajem i grupa; dopasowanie idzie takze SRODKIEM nazwy, bo przy przedrostkach zrodla szukanie od poczatku byloby bezuzyteczne (rozdz. 7.2)

**MessageType.CommandAgentConnectorList** — Odpowiednik `agent.plugin.list` dla drugiego z dwoch bytow

**MessageType.CommandAgentVersionGet** — `AgentVersion` niesie etykiete, opis zmiany, sprawce i sume kontrolna, ale NIE niesie tresci — bez tej komendy podgladu wersji nie ma z czego zlozyc, a przywrocenie jest jedynym sposobem zobaczenia, co w wersji stalo

**MessageType.CommandAgentAssignmentList** — `workspace.agent.assign` zapisuje przynaleznosc, ale zadna komenda nie odczytuje jej OD STRONY EKSPERTA, wiec biblioteka nie ma skad wziac liczby

**MessageType.CommandAgentModulesSet** — LISTA PUSTA `[]` ZNACZY BRAK OGRANICZENIA, czyli dostepnosc wszedzie, i jest to stan wyjsciowy; pominiecie pola znaczy „bez zmiany". Ta sama zasada co przy `memoryLevels`, gdzie zbior pusty tez niesie znaczenie wlasne

**MessageType.CommandAgentIsolationSet** — Wartosc obowiazuje wszedzie, gdzie ekspert dziala, dopoki nie nadpisze jej regula zapisana na poziomie zasiegu bardziej szczegolowym

**MessageType.CommandAgentPolicyGet** — Transparentnosc, nie bramka — komenda niczego nie zapisuje i niczego nie rozstrzyga

**MessageType.CommandAgentSubagentSet** — Granica pietnastu jest parametrem technicznym platformy, nie decyzja produktowa — ta sama, ktora zna `subagent.spawn`

**MessageType.CommandAgentPermissionRemove** — `agent.permission.set` wylacznie USTAWIA wartosc wpisu, wiec po pierwszym zawezeniu wiersz zakresu zostaje w wykazie na zawsze, choc jego wartosc wrocila do przyznanej

**MessageType.CommandChannelCheck** — Narzedzie pomocnicze, nie bramka: wynik nie warunkuje zapisu eksperta ani wyslania tury

**MessageType.CommandChannelCredentialStatus** — NIE ODDAJE TRESCI POSWIADCZENIA I ODDAWAC JEJ NIE MOZE

**MessageType.CommandExtensionSearch** — Dopelnia extension.list, ktory zawezal wylacznie rodzajem i stanem zainstalowania; dopasowanie idzie takze SRODKIEM nazwy, bo pozycje niosa przedrostki zrodla

**MessageType.CommandExtensionDetailGet** — Extension niesie sam naglowek pozycji, wiec karta szczegolow nie ma dzis z czego powstac

**MessageType.CommandExtensionCollectionApply** — Wynik jest bilansem, nie potwierdzeniem: pozycja odrzucona wraca z powodem, zamiast znikac po cichu

**MessageType.CommandExtensionRegistryList** — Wyliczenie ExtensionOrigin ma dzis dwie wartosci, wiec rejestr organizacji nie ma jak sie w katalogu pokazac

**MessageType.CommandExtensionUpdateCheck** — Extension niesie wlasna wersje, ale nie wersje dostepna, wiec wskaznika aktualizacji nie ma dzis z czego zlozyc

**MessageType.CommandExtensionVersionPin** — Wersja jest dzis polem opisowym, wiec nie ma czego przypiac

**MessageType.CommandExtensionVersionRollback** — Wykonuje sie od razu; ustawienie extension.rollback.confirm wlacza potwierdzenie, ktore nie warunkuje wykonania

**MessageType.CommandExtensionBundleInstall** — Wynik jest bilansem: pozycja odrzucona wraca z powodem

**MessageType.CommandExtensionHistoryList** — Rejestr niesie stan biezacy, nie droge, ktora do niego doprowadzila

**MessageType.CommandExtensionAdminBulk** — Wynik jest bilansem przyjetych i odrzuconych, nie pojedynczym potwierdzeniem

**MessageType.CommandExtensionToolList** — Odpowiedz na tools/list, resources/list i prompts/list nie ma dzis w kontrakcie ani komendy, ani struktury, w ktora mialaby wejsc

**MessageType.CommandExtensionToolCall** — Wywolanie jest probne: idzie poza kontekstem eksperta i niczego mu nie przypisuje

**MessageType.CommandExtensionTransportSet** — Transport da sie dzis wpisac wylacznie do nieprzezroczystej konfiguracji, wiec rdzen i okno moga rozumiec to pole inaczej

**MessageType.CommandExtensionCredentialBind** — Moduł operuje na referencjach; tresc poswiadczenia nigdy nie opuszcza rdzenia i nie wchodzi w to zadanie

**MessageType.CommandExtensionSecretList** — Oddaje wylacznie odwolania i metryke; tresc sekretu nie opuszcza rdzenia

**MessageType.CommandExtensionPermissionList** — Extension niesie nieprzezroczysta konfiguracje, ktora deklaracja uprawnien nie jest

**MessageType.CommandExtensionPermissionGrant** — Nadanie jest jedyna kontrola; brak nadania nie wstrzymuje instalacji ani wlaczenia (zasada zero blokad)

**MessageType.CommandExtensionManifestScan** — Wynik jest SYGNALEM, nie brama: nie wstrzymuje instalacji ani wlaczenia

**MessageType.CommandAppsStageList** — Odpowiednik odczytu dla zdarzenia apps.build.changed — bez niego wykaz etapow zaczyna pusty przy kazdym wejsciu w modul, choc rdzen trzyma etapy w bazie

**MessageType.CommandAppsStageSave** — AppStage.ownerAgentId jedzie dzis wylacznie w kierunku od rdzenia, wiec przypisania wykonawcy nie ma jak zlecic

**MessageType.CommandAppsProductLinkList** — Powiazania nie sa czynne domyslnie — kazde jest swiadoma decyzja Operatora

**MessageType.CommandAppsArchitectureValidate** — Dzis pole validationIssues jest zapisywane i oddawane bez wyliczenia, wiec jego pustka nie orzeka o poprawnosci ukladu. Zastrzezenia sa OSTRZEZENIEM, nie brama — nie wstrzymuja pracy w warsztatach

**MessageType.CommandAppsArchitectureVersionList** — AppArchitecture.version jest liczba, ale historii kontrakt nie niesie, wiec porownania nie ma z czym zestawic

**MessageType.CommandAppsPreviewStart** — Podglad na zywo jest narzedziem warstwy 1 Frontend Workspace, a rdzen nie stawia dzis serwera produktu

**MessageType.CommandAppsEndpointList** — AppComponent niesie apiContract jako JEDEN tekst, wiec punkt koncowy nie jest dzis bytem, ktory dalby sie wyliczyc

**MessageType.CommandAppsEnvironmentList** — Wyliczenie AppDeployEnvironment ma trzy wartosci stale, a opracowanie zada dowolnej liczby srodowisk konfigurowanych przez Operatora

**MessageType.CommandAppsEnvironmentVariableList** — Wartosc sekretna oddawana jest jako referencja, nigdy jako tresc

**MessageType.CommandAppsServiceLogRead** — Podglad na zywo idzie zdarzeniem apps.service.log

**MessageType.CommandAppsArtifactList** — AppDeployment niesie wersje i odnosnik logu, ale artefakt nie jest dzis bytem umowy

**MessageType.CommandAppsDeploymentLogRead** — AppDeployment niesie logRef — sam odnosnik, bez tresci — wiec podgladu logu nie ma dzis czym wypelnic. Podglad na zywo idzie zdarzeniem apps.deployment.log

**MessageType.CommandAppsPackageValidate** — Zastrzezenia sa OSTRZEZENIEM, nie brama: nie wstrzymuja publikacji

**MessageType.CommandAppsPackageSign** — Klucz wydawcy lezy w warstwie sekretow i nie wchodzi w to zadanie ani w odpowiedz

**MessageType.CommandSpeechAudioUpload** — Jest to brak platformowy, nie brak jednego okna: speech.transcribe bierze audioRef, czyli SCIEZKE PLIKU na maszynie silnika, a nagranie z mikrofonu istnieje wylacznie jako bajty w pamieci karty. Ten sam brak wstrzymuje dyktowanie w oknie komunikacji (client/src/okno-komunikacji/dyktowanie/dostarczenie-nagrania.ts) i mikrofon Voice Console — jedna komenda otwiera obie drogi naraz. Dzwiek nie opuszcza maszyny rdzenia: bajty ida do magazynu nagran, nie do sieci

**MessageType.CommandSpeechAudioFetch** — Domyka pare do speech.audio.upload i do translate.speech.synthesize: dzis rdzen oddaje odnosniki (AssistantActivityEntry.audioRef, AssistantVoiceCommandResponse.speechRef), lecz nie ma komendy, ktora by po nie siegnela — odslucha nie ma czym wykonac

**MessageType.CommandSpeechWakeGet** — Brak modelu frazy jest ODPOWIEDZIA, nie awaria — tak samo jak brak silnika w speech.availability.get

**MessageType.CommandSpeechWakeSet** — Pola pominiete zostaja bez zmian

**MessageType.CommandSpeechListenStart** — Rdzen slucha strumienia i oglasza zdarzeniami: czesciowa transkrypcje oraz wykrycie frazy wybudzajacej. Nasluch nie jest bramka — Operator zatrzymuje go speech.listen.stop, a rdzen nie zatrzymuje go sam

**MessageType.CommandSpeechListenStop** — Zatrzymanie nasluchu, ktorego nie ma, nie jest bledem

**MessageType.CommandAssistantActivityFlag** — Dzis wyroznienie wpisu nie mialoby gdzie zamieszkac: AssistantActivityEntry nie niesie takiego pola, wiec okno pokazywaloby wyroznienie, ktorego rdzen nie pamieta

**MessageType.CommandMemoryRetentionGet** — Odczyt do pary z memory.retention.set

**MessageType.CommandMemoryRetentionSet** — Zasada obejmuje zapisy kolejne, a wpisow zastanych nie rusza wstecz — inaczej zmiana reguly kasowalaby ustalenia, na ktore Operator sie nie umawial

**MessageType.CommandMemoryContextList** — Dzis memory.toggle przestawia poziomy zasiegu i to jedyny kontekst przelaczalny; nazwy zestawu nie ma gdzie odlozyc

**MessageType.CommandMemoryContextSave** — Puste contextId zaklada nowy, podane zmienia istniejacy — tak samo jak memory.set

**MessageType.CommandMemoryContextActivate** — Aktywacja nie kasuje kontekstu poprzedniego — wraca sie do niego tym samym wywolaniem

**MessageType.CommandMemoryContextDelete** — Wpisow pamieci nie kasuje — kontekst jest zestawem wskazan, a nie ustalenie budowy tresci

**MessageType.CommandToolsScopeList** — Dzis takiego zakresu nie ma dla profilu: agent.connector.* i agent.permission.set dotycza eksperta modulu Agents i nie siegaja profilu asystenta

**MessageType.CommandToolsScopeSet** — Zakres jest nastawa zasiegu, nie bramka: stanem wyjsciowym jest pelny dostep bez limitu, zgodnie z zasada zero blokad

**MessageType.CommandClipboardList** — Schowek jest bytem klienta, ktorego rdzen dzis nie zna, wiec historia ginie wraz z karta — komenda daje jej trwalosc

**MessageType.CommandClipboardPush** — Powtorzenie tresci identycznej nie mnozy wpisow — podnosi wpis zastany na czolo wykazu

**MessageType.CommandClipboardPin** — Przypiety nie wygasa wraz z zasada retencji historii

**MessageType.CommandSnippetList** — Skrot dziala we wszystkich polach tekstowych platformy, wiec slownik nalezy do rdzenia, nie do jednego okna

**MessageType.CommandSnippetSet** — Puste snippetId zaklada nowy, podane zmienia istniejacy

**MessageType.CommandLauncherHotkeyGet** — Skrot globalny rejestruje powloka programu okiennego, wiec brak wsparcia powloki jest ODPOWIEDZIA, nie awaria — klient nie obiecuje wtedy skrotu, ktory nikogo nie obudzi

**MessageType.CommandLauncherHotkeySet** — Skrot zajety przez inny program nie jest bledem zapisu — zapis zostaje, a odpowiedz mowi, ze rejestracja sie nie udala

**MessageType.CommandContextUsageGet** — ZALEZNOSC: wymaga tokenizatora jako podsystemu rdzenia — bez niego liczba tokenow bylaby wartoscia wzieta znikad. Do czasu jego zbudowania komenda ma zwracac available falsz wraz z powodem, tak jak speech.availability.get przy braku silnika

**MessageType.CommandDesignAssetContentGet** — Pole uri zasobu jest sciezka w systemie plikow rdzenia (magazyn oddaje filepath), wiec przegladarka nie wczyta spod niego niczego — take wtedy, gdy zasob powstal bez zarzutu. Komenda dotyczy KAZDEGO zasobu magazynu, nie tylko obrazu: jeden magazyn obsluguje rodziny design.*, document.*, media.* i archive.*. Rdzen odmawia zasobu, ktorego nie zna, i zasobu, ktorego tresci nie ma pod suma kontrolna — nigdy nie oddaje bajtow zastepczych

**MessageType.CommandDesignAssetExport** — Rozni sie od image.convert celem: konwersja zaklada NOWY zasob w magazynie i tam sie konczy, eksport oddaje bajty gotowe do zapisania poza produktem

**MessageType.CommandDesignAssetExportBatch** — Pokrywa eksport zbiorczy Assets Panel oraz zestawy rozmiarow kampanii; odmowa jednego zasobu NIE wstrzymuje pozostalych — wynik niesie bilans przyjetych i odrzuconych

**MessageType.CommandDesignCollectionCreate** — Etykiety zasobu juz sa (design.asset.tag.set), ale kolekcja jest bytem osobnym: ma nazwe, opis i porzadek, a etykieta jest tylko slowem

**MessageType.CommandDesignCollectionAssign** — Zestaw jest DOKLADKA albo ODJECIEM, nie zastapieniem — inaczej niz przy etykietach, bo kolekcja bywa duza i przepisywanie jej w calosci przy kazdej zmianie jest droga do zgubienia zawartosci

**MessageType.CommandDesignCollectionList** — Bez tego kolekcja zalozona nie mialaby jak wrocic na ekran po odswiezeniu

**MessageType.CommandDesignPromptTemplateSave** — Dzis historia promptow i szablony zyja w oknie do zamkniecia karty przegladarki i tyle o nich wiadomo

**MessageType.CommandDesignPromptHistoryList** — Domyka tez prowenancje: dzis zasob niesie pole promptId, ktorego rdzen NIE wypelnia, bo nie ma przekladu klucza wiersza promptu na kod kontraktu — ta komenda ten przeklad wnosi

**MessageType.CommandDesignBoardVersionRestore** — Przywrocenie ZAKLADA nowa wersje z ukladu sprzed przywrocenia, zeby cofniecie sie samo nie kasowalo stanu, ktory Operator wlasnie porzucil

**MessageType.CommandDesignBoardExport** — Dzis kompozycja jezdzi do rdzenia i z powrotem jako uklad warstw i nie ma drogi wyjscia poza rdzen

**MessageType.CommandDesignAnnotationSet** — Pole note warstwy niesie JEDNO zdanie bez autora i bez watku; opracowanie wymaga watkow i oznaczen osob, a tego jedno pole nie unosi

**MessageType.CommandDesignPresenceReport** — Rdzen rozglasza je pozostalym zdarzeniem design.board.presence. Zgloszenie jest ULOTNE — nie zapisuje sie w bazie, bo polozenie kursora sprzed godziny nie jest wiedza o niczym

**MessageType.CommandDesignTokensetSave** — Dzis Tokens & System Panel czyta zetony z motywu obowiazujacego i nie ma ich gdzie odlozyc — kontrakt nie zna bytu zestawu zetonow

**MessageType.CommandDesignTokensetExport** — Klient sklada dzis zmienne CSS, SCSS, konfiguracje Tailwind i modul JavaScript sam i nie potrzebuje do tego rdzenia; ta komenda jest potrzebna dla postaci, ktorych przegladarka zlozyc nie moze, oraz dla wydania idacego DO INNEGO MODULU zamiast do pliku

**MessageType.CommandDesignTokensetImport** — Rdzen NIE nadpisuje motywu produktu — motyw jest wlasnoscia powloki; import zaklada byt obok niego i oddaje roznice wobec zetonow wskazanego motywu

**MessageType.CommandDesignStyleguidePublish** — Przegladarka sklada przewodnik sama i oddaje go plikiem, ale wydania go do Library albo Studio nie ma czym zlecic

**MessageType.CommandImageVectorize** — Rodzina image.* zna przeksztalcenie geometryczne, poprawke, konwersje formatu, powiekszenie i wyciecie tla — zamiany rastra na sciezki nie zna zadna z nich, a jest to funkcja opracowania modulu

**MessageType.CommandImageLayersSplit** — Zasila rozdzielenie generacji na warstwy oraz zaznaczanie obiektu; bez silnika segmentacji ODMAWIA, nazywajac brak — nigdy nie oddaje calego obrazu jako jednej warstwy udajac rozklad

**MessageType.CommandImageCompose** — Pokrywa znak wodny, branding wsadowy, osadzenie w ramce urzadzenia i warstwy rastrowe — cztery funkcje opracowania, ktore wszystkie potrzebuja jednego: zlozenia obrazu na obrazie

**MessageType.CommandRoundtableDebateGet** — Obszar nie ma dzis zadnej komendy odczytu, wiec okno otwarte w trakcie debaty zna wylacznie to, co uslyszalo zdarzeniem od swojego otwarcia

**MessageType.CommandStudioDocumentFormatSet** — Dzis format nadaje rdzen przy wczytaniu i zadna komenda go nie zmienia: studio.document.save przyjmuje documentId, content, title i createVersion, ale nie format

**MessageType.CommandStudioCommentResolve** — Watek nie znika: rozwiazanie jest stanem, nie usunieciem

**MessageType.CommandStudioTrackingSet** — Przy wlaczonym sledzeniu kazdy zapis odklada wstawienia i usuniecia jako zmiany do decyzji, zamiast nadpisywac tresc

**MessageType.CommandStudioTrackingDecide** — Decyzja zapisuje sie w tresci dokumentu i zaklada wersje

**MessageType.CommandStudioAssetEmbed** — Realizuje powiazanie Design do Studio

**MessageType.CommandStudioOperationDelete** — Operacji fabrycznej nie usuwa — na nia odpowiada odmowa nazywajaca ten fakt

**MessageType.CommandStudioChainRun** — Przebieg prowadzi petla wykonawcza okna, a kazdy krok odklada wlasna propozycje zmiany

**MessageType.CommandStudioBatchRun** — Kazdy dokument dostaje wlasne zadanie petli; odmowa jednego nie wstrzymuje pozostalych

**MessageType.CommandStudioAnnotationAdd** — Dzis czynnosc ta idzie droga generyczna window.action i wraca odmowa not_found, bo katalog akcji nie ma jej wiersza

**MessageType.CommandStudioDiffReportExport** — Wynik jest zasobem magazynu rdzenia

**MessageType.CommandStudioDiffVisual** — Sluzy tam, gdzie roznica tekstowa nie widzi zmiany ukladu

**MessageType.CommandStudioSearchSemantic** — Uzupelnia wyszukiwanie wzorca w studio.diff.compare

**MessageType.CommandStudioDiffSource** — Sluzy kontroli, czy redakcja nie odeszla od zrodla

**MessageType.CommandStudioProposalDecide** — Dzis decyzja zapada wylacznie w kliencie i rdzen o niej nie wie, wiec propozycja zostaje w nim nierozstrzygnieta

**MessageType.CommandStudioBranchCreate** — Sluzy prowadzeniu dwoch redakcji obok siebie

**MessageType.CommandStudioBranchMerge** — Konflikt nierozstrzygniety wraca w wyniku zamiast byc rozstrzygniety domyslem

**MessageType.CommandStudioPreviewRender** — Zwraca strony jako zasoby, wiec podglad pokazuje UKLAD, a nie sam tekst

**MessageType.CommandStudioPdfSplit** — Kazda czesc jest osobnym zasobem

**MessageType.CommandStudioSecurityEncrypt** — Czynnosc jest jawnym, odwracalnym ustawieniem Operatora i niczego nie warunkuje

**MessageType.CommandStudioSecurityRedact** — Czynnosc jest nieodwracalna dla wyniku, dlatego zrodlo zostaje nietkniete

**MessageType.CommandStudioSecuritySensitiveDetect** — Czynnosc CZYTA i niczego nie zmienia

**MessageType.CommandStudioIngestQueueAdd** — Dzis kolejka jest wylacznie kliencka i ginie z odswiezeniem okna

**MessageType.CommandStudioIngestRecognize** — Komenda document.text.extract robi to samo waskim wejsciem — jeden jezyk, bez silnika, bez progu i bez ukladu

**MessageType.CommandStudioIngestUrl** — Migawke strony oddaje browser.snapshot.get, ale wymaga okna modulu Browser i nie prowadzi do dokumentu Studia

**MessageType.CommandStudioIngestDeviceList** — Bez tego wykazu wybor urzadzenia nie ma z czego powstac

**MessageType.CommandStudioStyleSave** — Zmiana stylu przestawia wszystkie miejsca dokumentu, ktore go uzywaja

**MessageType.CommandStudioStyleDelete** — Stylu fabrycznego nie usuwa — odpowiada odmowa nazywajaca powod

**MessageType.CommandStudioPageSetupSet** — Zmiana formatu przelicza uklad i oddaje bilans tego, co sie nie zmiescilo

**MessageType.CommandStudioTableStructureEdit** — Szerokosci kolumn zostaja policzone, nie zerowe

**MessageType.CommandStudioDocumentImportPdf** — Odzyskanie jest odtworzeniem, nie odczytem, wiec odpowiedz niesie bilans; PDF ze samych skanow kieruje na rozpoznanie tekstu

**MessageType.CommandStudioDocumentCopy** — Kopia jest osobnym dokumentem, nie drugim odwolaniem do tego samego

**MessageType.CommandStudioDocumentExportFormat** — Format ubozszy niz dokument oddaje wykaz cech pominietych

**MessageType.CommandStudioTemplateDelete** — Szablonu fabrycznego nie usuwa — odpowiada odmowa nazywajaca powod

**MessageType.CommandStudioLockAdd** — Blokada obowiazuje w rdzeniu, przed dotknieciem tresci

**MessageType.CommandStudioLockRemove** — Blokade zdejmuje WYLACZNIE Operator — czynnosc modelu wraca odmowa nazywajaca powod

**MessageType.CommandStudioJournalRevert** — Czynnosc bedaca podstawa pozniejszej odmawia i nazywa zaleznosc, zamiast zostawic dokument w stanie niespojnym

**MessageType.CommandStudioModelChangesRevert** — Nie jest to przywrocenie wersji sprzed, bo to skasowaloby prace Operatora

**MessageType.CommandStudioAutosaveRun** — Nieudany zapis wraca nazwany, nie przemilczany

**MessageType.CommandStudioViewSet** — Gdzie da sie zrobic dwojako i obie drogi maja sens, wybor nalezy do Operatora i jest jawnym, odwracalnym ustawieniem

**MessageType.CommandStudioMarkupAdd** — Znakowanie modelu jest podpisane jako model

**MessageType.CommandStudioMarkupTypeDelete** — Rodzaju fabrycznego nie usuwa — odpowiada odmowa nazywajaca powod

**MessageType.CommandStudioAgentsSettingsGet** — Oba narzedzia sa domyslnie wylaczone

**MessageType.CommandStudioAgentsSettingsSet** — Wartosci ida zasiegami rodziny config, nie osobnym magazynem

**MessageType.CommandStudioPlanCreate** — Sam rozklad niczego nie uruchamia

**MessageType.CommandStudioPlanRun** — Gdy nastawa petli jest wylaczona, wraca odmowa nazywajaca brak nastawy, a nie cisza

**MessageType.CommandStudioAgentsClaim** — Fragment zajety przez kogo innego wraca odmowa nazywajaca wykonawce i zakres

**MessageType.CommandStudioAgentsConflictsList** — Odlozone brzmienie nie przepada

**MessageType.CommandTerminalSessionClose** — Dzis zamkniecie karty zyje wylacznie w widoku klienta: powloka i jej procesy biegna dalej, a rdzen o zamknieciu nie wie

**MessageType.CommandTerminalSessionList** — Rdzen odtwarza karty przy starcie, ale klient po ponownym polaczeniu nie ma jak ich zobaczyc i zaczyna wykaz od pustego

**MessageType.CommandTerminalFileRead** — Bez tej komendy klient czyta manifest projektu poleceniem powloki, wiec wykrycie zadan zalezy od programu wypisujacego plik i od skladni kazdej z szesciu powlok

**MessageType.CommandTerminalProcessSuspend** — Dzis rdzen umie proces wylacznie zakonczyc, wiec dlugie zadanie da sie tylko ubic

**MessageType.CommandTerminalHostSave** — Bez tej komendy ksiazka zyje jedno posiedzenie przegladarki i ginie przy odswiezeniu strony

**MessageType.CommandTerminalHostRemove** — Karty juz otwarte do tego hosta biegna dalej

**MessageType.CommandTerminalScriptSave** — Bez tej komendy biblioteka zyje jedno posiedzenie, a jedyna droga jej zachowania jest wywoz do pliku

**MessageType.CommandTerminalScriptLint** — Analize prowadza programy spoza instalki Danaco Console, wiec odpowiedz mowi wprost, czy narzedzie bylo dostepne — pusty wykaz uwag przy braku narzedzia znaczylby falszywie tresc bez zastrzezen

**MessageType.CommandTerminalTunnelOpen** — Dzis tunel da sie zalozyc wylacznie poleceniem wydanym w karcie, a wtedy jego stan i przepustowosc sa niewidoczne

**MessageType.CommandTerminalKeyGenerate** — Haslo klucza wchodzi odwolaniem do sejfu, nigdy trescia — zgodnie z zasada zapisana w kontrakcie przy zmiennej srodowiska

**MessageType.CommandTerminalKeyImport** — Klucz wskazuje sie sciezka, a nie trescia: material kluczowy nie ma powodu przechodzic przez lacze

**MessageType.CommandTerminalKeyRemove** — Wpisy ksiazki hostow wskazujace ten klucz traca wskazanie i wracaja do klucza domyslnego konfiguracji maszyny

**MessageType.CommandTerminalWatchStart** — Kontrakt daje dzis wyzwalacz plikowy automatyce, a nie karcie powloki

**MessageType.CommandTerminalWatchStop** — Polecenie juz uruchomione biegnie dalej

**MessageType.CommandWorkspaceAgentUnassign** — Nie usuwa eksperta z biblioteki modulu Agents — znosi wylacznie jego przypisanie do tego projektu

**MessageType.CommandWorkspaceProjectStatusSet** — Droga do stanu 'paused', ktorego opracowanie wymaga, a ktorego zaden dzisiejszy uchwyt nie potrafi zapisac

**MessageType.CommandWorkspaceInstructionsVersionRestore** — Przywrocenie zaklada wersje nowa o tresci wersji wskazanej — historia nie jest przepisywana

**MessageType.CommandWorkspaceTaskUpdate** — Pola pominiete zostaja bez zmiany — wywolanie nie jest podmiana calego zadania

**MessageType.CommandWorkspaceTaskMove** — Osobno od workspace.task.update, bo przeciagniecie karty zmienia stan i porzadek naraz, a jest czynnoscia jednym ruchem myszy

**MessageType.CommandWorkspaceTaskDependencySet** — Rdzen odmawia zalozenia zaleznosci domykajacej cykl, bo cykl nie da sie ulozyc w czasie

**MessageType.CommandWorkspaceScheduleGet** — Zadanie bez obu granic czasu do harmonogramu nie wchodzi

**MessageType.CommandWorkspaceNoteSave** — Puste 'noteId' zaklada notatke nowa; podane zmienia istniejaca. Rdzen przy zapisie przelicza odnosniki tresci

**MessageType.CommandWorkspaceKnowledgeGraphGet** — Nazwa rodziny 'workspace.knowledge' nie miesza sie z rodzina 'knowledge.*': tamta prowadzi wskaznik ZNACZENIA, ta rysuje siec ODNOSNIKOW miedzy bytami projektu

**MessageType.CommandWorkspaceLibraryTextExtract** — Obejmuje rozpoznanie tekstu z obrazow i skanow

**MessageType.CommandWorkspaceLibraryDuplicateList** — Sam wykaz niczego nie scala

**MessageType.CommandWorkspaceLibraryDuplicateMerge** — Etykiety, kolekcje i wersje plikow scalanych przechodza na plik zachowany, zeby scalenie nie gubilo dorobku

**MessageType.CommandWorkspaceSearchProject** — Nie zastepuje library.file.search: tamta komenda przeszukuje biblioteke centralna, ta jeden projekt i wiecej niz pliki. Wyszukiwanie po ZNACZENIU prowadzi rodzina knowledge.

**MessageType.CommandLibraryMetadataGet** — Metadane osadzone czyta sie z bajtow, wiec ich odczyt jest kosztowny i wchodzi wylacznie na wyrazne zadanie

**MessageType.CommandLibraryMetadataSet** — Dzis formularz opisu nie ma dokad pojsc: library.tag.set zmienia wylacznie etykiety i kolekcje. Zapis rozglasza library.file.changed

**MessageType.CommandLibrarySchemaSet** — Pole o kodzie juz istniejacym jest zmieniane, nie dublowane

**MessageType.CommandLibraryTagList** — Bez tej komendy klient sklada slownik z etykiet plikow odczytanej strony wykazu, wiec pokazuje probke zamiast slownika

**MessageType.CommandLibraryTagUpdate** — Zmiana nazwy przechodzi na wszystkie zasoby noszace etykiete — inaczej powstalaby druga etykieta o tym samym znaczeniu

**MessageType.CommandLibraryTagMerge** — Zasoby noszace etykiete zrodlowa dostaja docelowa, a zrodlowa znika ze slownika

**MessageType.CommandLibraryTagRemove** — Usuniecie etykiety uzywanej zada potwierdzenia, bo zdejmuje ja z zasobow, ktorych zadanie nie wymienia

**MessageType.CommandLibraryCollectionList** — Bez tej komendy klient zna wylacznie identyfikatory kolekcji wyczytane z plikow, a nazwe kolekcji tylko w chwili jej zalozenia

**MessageType.CommandLibraryRuleSet** — Zapis uruchamia przeliczenie

**MessageType.CommandLibraryRuleRemove** — Zasoby przypisane przez regule kolekcji inteligentnej zostaja w kolekcji, ale przestaja byc przeliczane — inaczej usuniecie reguly oproznialoby kolekcje po cichu

**MessageType.CommandLibraryDuplicateScan** — Rozpoznanie przyblizone zada porownania calego zbioru, wiec jest czynnoscia rdzenia, nie okna

**MessageType.CommandLibraryDuplicateMerge** — Zasoby wchloniete trafiaja do archiwum, nie znikaja

**MessageType.CommandLibraryFixityCheck** — Klient nie ma czym tego zrobic, bo nie siega po bajty zasobu

**MessageType.CommandLibraryNameNormalize** — Przebieg probny pokazuje wynik bez zapisu

**MessageType.CommandLibraryAuditList** — Dziennik jest przyrostowy — wpisow nie da sie zmienic ani usunac ta rodzina komend

**MessageType.CommandLibraryRetentionSet** — Twarde usuniecie nie zachodzi samoczynnie: polityka najwyzej zglasza zasob do usuniecia

**MessageType.CommandLibraryFileMove** — Przeniesienie nie rusza tresci ani historii wersji

**MessageType.CommandLibraryFileArchive** — Czynnosc jest odwracalna komenda library.file.restore — to jest kosz repozytorium, nie usuniecie

**MessageType.CommandLibraryFileDelete** — Jedyna droga utraty danych repozytorium; bez potwierdzenia jest odmowa

**MessageType.CommandLibraryShareCreate** — Token jest jawny zgodnie z zasada jawnosci kluczy platformy

**MessageType.CommandLibraryShareRevoke** — Odnosnik przestaje dzialac, wpis zostaje w dzienniku audytu

**MessageType.CommandLibraryClassifyRun** — Wynik wraca sugestiami do przyjecia, nie zapisem — domyslnym zachowaniem modulu jest sugestia z akceptacja Operatora

**MessageType.CommandLibrarySuggestionApply** — Przyjecie wykonuje czynnosc, ktora sugestia opisuje; odrzucenie zdejmuje ja z wykazu

**MessageType.CommandLibraryDiffCompare** — Dokumenty binarne porownuje po tekscie z nich wydobytym i mowi o tym w odpowiedzi

**MessageType.CommandDesignVectorPathSet** — Wezly przychodza w calosci — zmiana jednego wezla idzie ta sama droga co narysowanie sciezki, zeby nie bylo dwoch prawd o jej ksztalcie

**MessageType.CommandDesignVectorShapeAdd** — Ksztalt powstaje od razu jako wezly, wiec da sie go dalej edytowac pioram — nie jest osobnym bytem, ktory potem trzeba zamieniac

**MessageType.CommandDesignVectorBoolean** — Sciezki zrodlowe znikaja albo zostaja wedle wskazania — bo suma dwoch ksztaltow bywa krokiem posrednim, a bywa wynikiem koncowym

**MessageType.CommandDesignVectorTextPath** — Zamiana w kontury jest nieodwracalna dla wyniku, dlatego tekst zrodlowy zostaje

**MessageType.CommandDesignVectorOptimize** — Odpowiedz niesie ubytek zmierzony, zeby Operator widzial, ile naprawde ubylo, zamiast czytac obietnice

**MessageType.CommandDesignVectorSymbolSet** — Zmiana definicji propaguje do wszystkich instancji — po to symbol jest

**MessageType.CommandDesignFrameSet** — Ramka jest ekranem makiety: to ona, a nie kanwa, wyznacza obszar wydania

**MessageType.CommandDesignFrameRemove** — Warstwy ramki zostaja na kanwie — usuniecie ramki nie jest usunieciem pracy, ktora w niej lezala

**MessageType.CommandDesignLayoutAuto** — Rdzen oddaje warstwy po przeliczeniu, wiec klient nie liczy ukladu drugi raz i nie ma jak sie rozjechac

**MessageType.CommandDesignFrameResizeApply** — Osobno od design.frame.set, bo tam zmiana rozmiaru jest zapisem nastawy, a tu przeliczeniem ukladu wedle wiezi

**MessageType.CommandDesignMockupGenerate** — Wynik jest makieta do poprawienia, nie projektem koncowym, i tak sie o nim mowi

**MessageType.CommandDesignMockupImport** — Odpowiedz mowi, czego rdzen nie rozpoznal — bilans zamiast ciszy

**MessageType.CommandDesignColorVisionSimulate** — Wynik jest nowym zasobem — zrodlo zostaje nietkniete

**MessageType.CommandDesignColorAccessibilityAudit** — Czynnosc CZYTA i niczego nie zmienia

**MessageType.CommandDesignStockSearch** — Brak skonfigurowanego dostawcy jest ODPOWIEDZIA nazywajaca brak, nie cisza

**MessageType.CommandDesignPrintPreflight** — Czynnosc CZYTA i niczego nie zmienia; zastrzezenie z waga blad zatrzyma druk, wiec Operator ma je zobaczyc TUTAJ, a nie w drukarni

**MessageType.CommandDesignLargeformatTile** — Odpowiedz niesie uklad kafli, wiec Operator wie, ktory kafel gdzie idzie, zanim cokolwiek wydrukuje

**MessageType.CommandDesignChartRender** — Dane przychodza seriami, a nie obrazem — wykres da sie wiec przerysowac po zmianie liczb, zamiast rysowac go od nowa

**MessageType.CommandDesignPhotoCrop** — Wynik jest wariantem zrodla — oryginal zostaje nietkniety

**MessageType.CommandDesignPhotoUpscale** — Gdy kanal modelu obrazowego stoi, liczy kanal; gdy nie stoi, liczy rachunek wkompilowany — odpowiedz mowi to polem computedBy

**MessageType.CommandDesignPhotoInpaint** — Gdy kanal modelu obrazowego stoi, liczy kanal; gdy nie stoi, liczy rachunek wkompilowany — odpowiedz mowi to polem computedBy

**MessageType.CommandDesignPhotoExpand** — Gdy kanal modelu obrazowego stoi, liczy kanal; gdy nie stoi, liczy rachunek wkompilowany — odpowiedz mowi to polem computedBy

**MessageType.CommandDesignPhotoBackgroundRemove** — Gdy kanal modelu obrazowego stoi, liczy kanal; gdy nie stoi, liczy rachunek wkompilowany — odpowiedz mowi to polem computedBy

**MessageType.CommandAodMuteGet** — Wyciszenie przeterminowane nie wchodzi do wykazu — Operator nie ma go odklikiwac

**MessageType.CommandAodMuteSet** — Zniesienie idzie ta sama komenda z polem muted rownym falszowi — odwracalnosc jednym ruchem jest wymogiem, nie wygoda. Odpowiedz oddaje wykaz wyciszen PO zmianie, zeby powloka nie musiala pytac drugi raz. Zmiana rozglasza sie zdarzeniem aod.mute.changed na pozostale powloki Operatora

**MessageType.CommandAodSignalReport** — Nosnik dla trzech klas rozdz. 3.2, ktorych rdzen nie widzi wlasna telemetria: wyniku kontroli jakosci, harmonogramu przebiegow automatyk i powtarzalnosci czynnosci Operatora. Odpowiedz mowi wprost, czy sygnal wpadl w wyciszenie i ktore — sygnal wyciszony odklada sie nadal, bo wyciszenie wstrzymuje UJAWNIENIE, a nie zapis

**MessageType.CommandAodSignalList** — Odpowiedz nazywa wyciszenia, ktore sygnaly wstrzymaly, i liczbe wstrzymanych — cisza, po ktorej Operator nie wie, ze cos milczy, jest gorsza od braku wyciszenia

**MessageType.CommandMemoryDisableList** — Odczyt do pary z memory.disable.set; okno konfiguracji buduje z niego zakres ustawien „Pamiec"

**MessageType.CommandMemoryDisableSet** — Zniesienie idzie ta sama komenda z polem disabled rownym falszowi: odwracalnosc jednym ruchem. WYLACZENIE NIE KASUJE TRESCI — tym rozni sie od memory.delete; wpis wylaczony zostaje na miejscu i wraca w calosci po zniesieniu. Odpowiedz oddaje wykaz wylaczen PO zmianie

**MessageType.CommandNotificationList** — Filtry zawezaja wykaz do klas, wag, stanow albo zrodla; pominiete zwracaja rejestr niezamkniety. Licznik zdarzen nowych idzie zawsze, niezaleznie od filtru — plakietka paska kontekstu liczy caly rejestr, nie widok

**MessageType.CommandNotificationAcknowledge** — Wykaz pusty oznacza wszystkie zdarzenia nowe, czyli dzialanie zbiorcze centrum

**MessageType.CommandNotificationResolve** — Dzialanie, ktore je zamyka, wykonuje rodzina wlasciwa — ta komenda zapisuje wylacznie skutek w rejestrze centrum

**MessageType.CommandNotificationSnooze** — Po niej zdarzenie wraca do stanu nowe i znow liczy sie do plakietki — odlozenie nie jest zamknieciem

**MessageType.EventAuthChanged** — Sekcja Uwierzytelnianie Okna Ustawien odswieza sie bez odpytywania

**MessageType.EventDeviceChanged** — Rozglaszane do WSZYSTKICH polaczonych urzadzen, zeby uniewaznienie jednego bylo natychmiast widoczne na pozostalych ekranach — tej samej drogi uzywa odzyskanie konta, ktore uniewaznia tokeny wydane przed zmiana hasla

**MessageType.EventAdvisorConsulted** — Zdarzenie istnieje dla JAWNOSCI: rada, ktorej nie widac w strumieniu, jest przeslanka, ktorej w aktach nie ma — a agent na niej dziala. Niesie okno, doradce, podstawe doboru i skrot rady; pelna tresc wraca wynikiem komendy advisor.consult i zapisem dziennika konsultacji

**MessageType.EventAppsWorkspaceChanged** — Bez tego zdarzenia apps.workspace.update byl jedyna komenda zapisu modulu, ktora nic nie rozglaszala — drugie okno nie dowiadywalo sie o zmianie

**MessageType.EventSessionToolAttached** — Skoro model dostaje narzedzie, ktorego nie mial, Operator MUSI to zobaczyc — kazde posuniecie widac na ekranie

**MessageType.EventSessionToolDetached** — Widocznosc obowiazuje w OBIE strony: zestaw, ktory urosl na oczach Operatora, nie moze skurczyc sie po cichu — bez tego zdarzenia sasiednie okno pokazywaloby narzedzie, ktorego model juz nie ma

**MessageType.EventAppsPreviewChanged** — Bez tego zdarzenia podglad na zywo bylby odswiezaniem na zadanie, czyli nie podgladem na zywo

**MessageType.EventAppsDeploymentLog** — Osobne od apps.build.changed, ktore niesie STAN wdrozenia, a nie jego dziennik

**MessageType.EventSpeechListenPartial** — Tekst jest nietrwaly i bywa poprawiany kolejnym zdarzeniem — trwaly zapis daje dopiero speech.transcribe

**MessageType.EventAutomationExecutionLogged** — Nazwa rozna od komendy automation.execution.log, bo kontrakt nie dopuszcza tej samej nazwy w obu wykazach

**MessageType.EventDesignBoardChanged** — Dzis kompozycja jezdzi w obie strony komendami, ale drugie okno i drugie polaczenie nie dowiaduja sie o zmianie niczym — zdarzenia obszaru jest dokladnie jedno i dotyczy zasobu

**MessageType.EventDesignBoardPresence** — Zdarzenie ULOTNE — nie ma odpowiednika w bazie i nie jest odtwarzane po ponownym polaczeniu

**MessageType.EventAlertTriggered** — Zdarzenie jest droga alertu do Always On Display i do powiadomienia w aplikacji; rejestr odpytywany komenda alert.trigger.list sam by tam nie dotarl, bo alert ma znaczenie wtedy, gdy Operator nie patrzy

**MessageType.EventAodMuteChanged** — Rozgloszenie sciga cisza pozostale powloki Operatora: bez niego Operator wyciszalby w jednej, a sugestie wchodzilyby w drugiej

**MessageType.EventMemoryDisableChanged** — Rozgloszenie idzie na pozostale powloki Operatora, zeby wykaz pamieci w kazdej z nich pokazywal ten sam stan wlaczenia

**MessageType.EventNotificationRaised** — Rozglaszane do wszystkich polaczen Operatora, zeby plakietka i kolumna centrum nie musialy odpytywac rdzenia w tle; to samo zdarzenie jest podstawa Toastu w chwili wystapienia

**MessageType.EventNotificationChanged** — Rozglaszane do wszystkich polaczen, zeby obsluga na jednym urzadzeniu byla natychmiast widoczna na pozostalych

**MessageType.EventConnectionUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventHomeUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventEnvironmentUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventModuleUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventWorkspaceUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventSessionUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventWindowUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventMessageUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventConfigUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventSettingsUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventAccessUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventAccountUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventIdentityUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventChannelUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventQueueUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventContextUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventActionUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventStreamUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventProgressUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventStudioUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventAutomationUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventBrowserUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventResearchUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventLibraryUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventTranslateUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventRoundtableUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventDesignUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventAssistantUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventTerminalUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventDeveloperUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventDiagnosticsUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventAppsUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventAgentUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventMemoryUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventOrchestrationUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventIsolationUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventComponentUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventExtensionUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventRoleUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventSubagentUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventAdvisorUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventMonitorUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventScheduleUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventModelUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventMobileUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventAodUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventAuthUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventSpeechUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventToolsUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventArchiveUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventDeviceUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventDocumentUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventHistoryUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventImageUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventKnowledgeUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventMailUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventMediaUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventPanelUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventRetentionUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventTeamUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventAlertUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventClipboardUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventHealthUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventLauncherUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventProvenanceUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventSnippetUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventUsageUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**MessageType.EventNotificationUnknown** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**Envelope** — Pola odpowiedzi i strumienia sa opcjonalne.

**Window.AgentId** — Agent nie zastepuje kanalu — kanal jest droga do modelu, agent tozsamoscia nalozona na te droge

**Window.CoordinatorWindowId** — Odpowiednik kolumny okno_komunikacji.okno_koordynatora_id

**Channel.AccountId** — Kanal mowi CZYM sie polaczyc, konto — CZYIM poswiadczeniem

**Environment** — Odpowiednik wiersza tabeli srodowisko

**Environment.Motto** — Brak motta jest stanem normalnym: puste pole ma zostac puste, nie wolno go zastapic tekstem ulozonym po stronie widoku

**Environment.ModuleCodes** — Puste dla srodowiska z panelem orkiestracji

**Module** — Odpowiednik wiersza tabeli modul

**Module.OperationalWindowCodes** — Okno rozmowy wystepuje w kazdym module

**Module.Kind** — Odpowiednik kolumny modul.rodzaj (migracja 076)

**Module.ConfiguredOnHome** — NIEZALEZNE OD `kind`: rodzaj mowi CZYM modul jest, to pole GDZIE sie go sklada

**ToolDeclaration** — Brak osobnego parsera intencji: polecenie wpisane w oknie rozmowy model wykonuje wywolaniem narzedzia

**SettingDefinition** — Niesie komplet metadanych potrzebnych do zbudowania pola formularza — klient nie zaszywa u siebie ani jednego klucza

**AccessPoint** — Punkt mowi, DO CZEGO model ma wglad; katalog roboczy modelu jest ustawieniem osobnym (klucz katalog.roboczy.podstawa)

**AccessGrant** — Okno ma ZBIOR nadan; kolejnosc i oznaczenie glownego maja znaczenie

**Account** — Poswiadczenia wchodza zadaniem i nie wychodza nigdy — odpowiedz niesie wylacznie znacznik hasCredential

**IdentityCategory** — Katalog obejmuje konstytucje, profil roli, ekspertyze, zasady bezpieczenstwa i zasady harnessu

**LoopState** — Bieg nie ma limitu obiegow — zamiast bramy licznikowej wchodzi przejrzystosc: licznik obiegow, obiegi bez postepu i jawny powod zatrzymania. Odpowiednik session.StanObiegu

**SessionPresence** — Zasila kontrolke powrotu na stronie glownej: rozlaczenie klienta nie konczy sesji ani procesow, wiec strona glowna musi wiedziec, ze sesja zyje i gdzie zyje. session.bind jest akcja powrotu wywolywana z tej kontrolki, nie ruchem otwierajacym

**SessionConfig** — NIE jest zbiorem flag programu: adapter dostawcy dopiero tlumaczy go na wiersz argumentow, zmienne srodowiskowe, plik ustawien, plik serwerow MCP, plik pamieci, pliki promptu i katalog roboczy. Ten sam model obsluguje kanal CLI, API, SDK, SSH i model lokalny — zmiana dostawcy zmienia adapter, nie model. Osiemnascie obszarow najwyzszego poziomu; warstwy tozsamosci, profilu roli i ekspertyzy sa polami obszaru systemPrompt, nie osobnymi obszarami. Kazdy obszar jest opcjonalny: brak obszaru znaczy wartosc odziedziczona z poziomu szerszego albo domyslna

**SessionConfigAccount** — Zrodlem prawdy dla konta i poswiadczenia jest rejestr kont (account.*) — konfiguracja sesji odwoluje sie do niego identyfikatorem i nigdy nie niesie tresci poswiadczenia

**SessionConfigProvider** — Zrodlem prawdy dla parametrow polaczenia jest rejestr kanalow modelu (channel.*); pola ponizej sa nadpisaniem per sesja, nie druga kopia rejestru

**SessionConfigSystemPrompt** — Tozsamosc jest ZAMIENIANA, nie dolaczana — domyslnym trybem jest ZASTAP. Warstwy odpowiadaja injection.WarstwaKonstytucja, injection.WarstwaProfil i injection.WarstwaEkspertyza; obszar STERUJE istniejacym silnikiem nakladki, nie powtarza go

**SessionConfigSystemPrompt.Mode** — Puste znaczy ZASTAP

**SystemPromptLayer** — Edycja tresci idzie WYLACZNIE komenda identity.document.set — konfiguracja sesji wskazuje kategorie, nie zastepuje katalogu tozsamosci

**SessionConfigTools** — Wykaz narzedzi platformy powstaje z kontraktu — tutaj rozstrzyga sie wylacznie ich dopuszczenie

**SessionConfigPermissions** — Pola bypassAllPrompts i promptHandlerTool sa punktami zastrzezonymi Operatorowi — model nie rozszerza wlasnych uprawnien

**SessionConfigMcp** — Zrodlem prawdy dla adresu i poswiadczenia mostu jest rejestr punktow dostepu (access.*) — wiazanie odwoluje sie do punktu identyfikatorem

**SessionConfigMemory** — Adapter zapisuje tresc do pliku pamieci katalogu roboczego — pole niesie tresc, nie sciezke pliku dostawcy

**SessionConfigConversationContext** — Przenoszeniem kompletu miedzy modulami steruje context.transfer

**SessionConfigWorkingDirectory** — Katalog roboczy jest ustawieniem OSOBNYM od dostepu i od srodowiska platformy. Rozejscie miedzy katalogiem ustawionym a faktycznie uzywanym pokazuje WorkingDirectoryResolution w konfiguracji obowiazujacej

**SessionConfigAdditionalDirectories** — Zrodlem prawdy dla zakresu wgladu jest rejestr nadan dostepu okna (access.grant.*) — dostep to nie katalog roboczy

**SessionConfigEnvironment** — To NIE jest srodowisko platformy — profil widocznosci modulow (TalkIn, WorkSpace, CodeStudio, MultitaskingAI) zyje w obszarze komend environment i pozostaje nietkniety

**EnvironmentVariable** — Wartosc tajna wchodzi odwolaniem do sejfu, nigdy trescia

**SessionConfigInputOutput** — Pola sa dziedzinowe, nie sa nazwami przelacznikow — adapter dopiero tlumaczy je na postac dostawcy

**SessionConfigRuntime** — Uruchomieniem steruje injection.Rozruch — jedyny spawner platformy

**SessionConfigSessionLifecycle** — Bieg naprawczy nie ma bramy licznikowej — zamiast niej stoi prog braku postepu i jawny powod zatrzymania; stan biegu niesie LoopState

**WorkingDirectoryResolution** — Degradacja ma byc widoczna w oknie, nie milczaca; sama degradacja nie wstrzymuje uruchomienia

**SessionConfigFieldCapability** — Pole bez odpowiednika u dostawcy jest wypisane JAWNIE wraz z powodem, zamiast byc pomijane po cichu — tak jak turnLimit i promptHandlerTool, ktorych zainstalowany program CLI nie zna. Zmiana wersji programu albo przejscie na API zmienia wylacznie adapter i te deklaracje

**AdapterCapabilities** — Klient pyta o nia przed zbudowaniem okna konfiguracji i wie, ktore pole dostawca zignoruje, zanim Operator je wypelni

**AgentVersion.Mode** — Pole NIEWYMAGANE i to niesie trzeci stan: brak wartosci oznacza prompt globalny, DOLACZ prompt dopisywany, ZASTAP odstepstwo jawne

**DiagnosticError.ErrorCode** — Odmowa komendy jest FAKTEM o niej samej, wiec kod nalezy do bledu, nie do jego kontekstu; kolumna `kod_bledu` trzyma go w bazie od migracji 043

**AgentLayer** — Warstwy eksperta DOPISUJA sie do promptu systemowego jako zakres uzytkownika i NIGDY go nie zastepuja — dlatego warstwa nie ma pola trybu. Prompt systemowy jest paczka built-in, obowiazujaca zawsze i jako pierwsza (rozstrzygniecie ustalenie budowy)

**Agent.Mode** — Pominiete znaczy DOLACZ: warstwy eksperta DOPISUJA sie do globalnego promptu systemowego z okna konfiguracji i ustawien. ZASTAP jest ODSTEPSTWEM od ustawien domyslnych, oznaczanym przez Operatora swiadomie — instrukcja eksperta STAJE SIE promptem systemowym, a globalny przestaje obowiazywac dla okien tego eksperta

**Agent.MemoryLevels** — Pole oddawane ZAWSZE, takze puste: LISTA PUSTA ZNACZY PAMIEC WYLACZONA W CALOSCI i jest to jedyny zapis wylaczenia ( MemoryLevel). Pominiecie pola zostawialoby Operatora z domyslem zamiast z odpowiedzia. Wartosc wyjsciowa nowego eksperta to cztery poziomy — stanem wyjsciowym platformy jest pelny dostep

**Agent.ModuleCodes** — Pole oddawane ZAWSZE, takze puste: LISTA PUSTA ZNACZY BRAK OGRANICZENIA, czyli dostepnosc we wszystkich modulach, i jest to stan wyjsciowy. Pominiecie pola zostawialoby Operatora z domyslem zamiast z odpowiedzia — ta sama zasada co przy `memoryLevels`

**Agent.SubagentLimit** — Wartosc wyjsciowa 15 — maksimum techniczne platformy. Zero znaczy Subagent Network wylaczony i jest to jedyny zapis wylaczenia, wiec `agent.subagent.set` z `enabled: false` zapisuje tu zero

**Extension.Origin** — Rozstrzyga WYLACZNIE stan wyjsciowy przy rejestracji; poza tym jest faktem do pokazania Operatorowi

**Subagent** — UWAGA: czym podagent JEST — oknem, procesem czy pozycja kolejki — nie zostalo rozstrzygniete przez ustalenie budowy

**AodStatus.ModuleId** — Do tej pozycji nakladka dochodzila modulu okreznie — przez window.list — i sugestia, ktorej okna rdzen nie znal, w wyciszenie modulu nie wpadala

**AodSuggestion.ModuleId** — Bez niego wyciszenie biezacego modulu nie ma po czym rozpoznac swojej sugestii

**AodSuggestion.EventClass** — Podstawa wyciszenia klasy zdarzen; puste znaczy klase nierozpoznana, a sugestia bez klasy w wyciszenie klasy nie wpada

**AuthMethod** — Sekret nie trafia do bazy — zostaje w niej wylacznie odwolanie; odpowiedz nie niesie ani skrotu hasla, ani skrotu PIN-u, ani materialu klucza prywatnego. PIN i Windows Hello sa WLASCIWE URZADZENIU: klucz nie opuszcza jego TPM, wiec na kazdej maszynie zaklada sie je osobno

**AuthSession** — Transportem jest gniazdo WebSocket, nie seria zadan HTTP, wiec ciasteczka nie ma — token wraca z auth.login i idzie w connection.hello. Wygasanie przesuwne, przedluzane aktywnoscia albo jawnie komenda auth.token.refresh

**Device** — Konto jest jedno, urzadzen dowolnie wiele; kazde niesie wlasny token dostepu, ktory Operator moze uniewaznic osobno

**MailAccount** — Poswiadczenia leza w sejfie rdzenia, nigdy w bazie ani w kontrakcie

**MailMessage** — Tresc jest wypelniona wylacznie przy mail.message.get — wykaz niesie same naglowki, bo ciagniecie tresci calej skrzynki byloby pobraniem archiwum przy kazdym pytaniu

**KnowledgeImageHit** — Osobna struktura od KnowledgeHit, bo obraz nie ma fragmentu tekstu do zacytowania — wskazaniem jest sam plik

**SessionTool** — Zyje W STANIE SESJI, nie w definicji eksperta: definicja pozostaje nietknieta, dolozenie przezywa rozlaczenie klienta i konczy sie wraz z sesja albo z usunieciem rozmowy

**ToolCatalogEntry** — Trzy pola nazewnicze nie sa ozdoba: nazwa pelna niesie przedrostek zrodla, nazwa skrocona jest tym, co Operator wpisuje, a opis ratuje przed cicha degradacja — nazwa :design-audit sama nie mowi nic (rozdz. 7)

**AgentVersionSnapshot** — Byt odrebny od `Agent`, bo wersja nie niesie ani identyfikatora eksperta, ani jego przypisan, ani stanu czynnosci — utrwalone jest to, co `agent.update` przyjmuje, i tylko to

**ChannelCheckResult** — Struktura wydzielona, bo ten sam ksztalt oddaje sprawdzenie wykonane recznie z okna i sprawdzenie wykonane przy podlaczeniu rozszerzenia rozdz. 13, punkt sterowania „Test polaczenia konektora")

**ChannelCredentialStatus** — Struktura celowo nie ma pola na wartosc: brak pola jest mocniejszym zabezpieczeniem niz zasada, ze pola sie nie wypelnia

**AppDependency** — AppComponent.dependsOn niesie same identyfikatory, wiec strzalki przeplywu danych nie ma dzis z czego narysowac

**ModelCallTrace** — Wiersz tabeli prowenancja_wywolanie

**ModelCallTrace.ProcessId** — To jest pole korelacyjne wspolne z LogEntry.processId i MonitorStatus.processId — po nim wpis dziennika laczy sie ze sladem

**ModelCallSpan** — Wiersz tabeli prowenancja_odcinek

**UsageAggregate** — Agregat powstaje z odczytu tabeli prowenancja_wywolanie, nie z osobnej tabeli zuzycia — dwie tabele o tych samych liczbach rozjechalyby sie przy pierwszej korekcie cennika

**UsageAggregate.CostWithoutPrice** — Wieksze od zera znaczy, ze koszt jest niepelny i okno ma to powiedziec

**AlertRule** — Wiersz tabeli alert_regula wraz z kanalami z tabeli alert_regula_kanal

**AlertTrigger** — Wiersz tabeli alert_wyzwolenie

**AlertTrigger.DeliveredChannels** — Roznica wobec regul.channels znaczy nieudane dostarczenie, nie brak wyzwolenia

**HealthProbe** — Wiersz tabeli kondycja_sonda

**HealthProbeResult** — Wiersz tabeli kondycja_wynik. Bez tej serii nie ma z czego policzyc ani procentu dostepnosci, ani budzetu bledow — dlatego wyniki sa wierszami, a nie ostatnim stanem nadpisywanym w miejscu

**HealthUptime.Samples** — Zero znaczy brak pomiaru, a nie dostepnosc zerowa — okno ma odroznic jedno od drugiego

**LibraryMetadata** — Piec pol profilu podstawowego (title, creator, subject, date, rights) wypelnia okno domyslnie; pozostale dziesiec nalezy do standardu Dublin Core Elements 1.1 i jest oplacalne raz, a nie przy kazdym rozszerzeniu profilu

**LibraryMetadata.Identifier** — DOI albo ISBN

**LibraryTechnicalMetadata** — Pola geolokalizacji sa jedynym zrodlem widoku mapy w Library Explorer — bez nich widok nie ma czego naniesc

**LibraryTag** — Dzis etykieta jest w kontrakcie samym napisem przy pliku i nie ma tozsamosci — bez niej nie da sie ani nadac jej barwy, ani polaczyc dwoch etykiet w jedna

**LibraryDuplicateGroup** — Podstawa rozpoznania wchodzi do wyniku, bo duplikat pewny i duplikat przyblizony zadaja innej decyzji Operatora

**LibraryFixityResult** — Suma oczekiwana i wyliczona stoja osobno, bo rownosc jest wnioskiem, a nie danymi

**LibraryStats** — Klient liczy dzis te wartosci z odczytanej strony wykazu, wiec orzeka o probce, a nie o repozytorium

**TerminalSshKey** — Czesc tajna nie opuszcza maszyny rdzenia i nie ma tu pola

**WorkspaceTask.Rank** — Napis, nie liczba: wstawienie karty miedzy dwie sasiednie nie ma wtedy przepisywac calej kolumny

**WorkspaceBacklink.TargetNoteId** — Puste znaczy odnosnik do strony jeszcze niezalozonej — to stan poprawny, nie usterka

**WorkspaceKnowledgeGraph.Truncated** — Bez tego pola przyciety graf wygladalby na kompletny obraz projektu

**WorkspaceCanvas.Scene** — Rdzen sceny nie rozbiera na wiersze: jej ksztalt nalezy do widoku, ktory ja rysuje

**StudioParagraphFormat.ListRestart** — Bez tego pola wznowienie rozdzielalo liste na dwie definicje, a zmiana jednej definicji przestawala przestawiac oba jej odcinki

**StudioFragmentLock.TemplateId** — Blokada wychodzaca z fromTemplate bez tego pola nie mowi, KTOREGO szablonu jest wzorcem, wiec Operator nie ma jak dojsc do jej zrodla

**StudioActor** — Pojecia agenta nie zaklada od nowa: agentId odpowiada Agent.id, subagentId odpowiada Subagent.id

**StudioAgentSettings** — Oba narzedzia sa DOMYSLNIE WYLACZONE i wlaczane jawnym, odwracalnym ustawieniem Operatora; wartosci ida zasiegami rodziny config, nie osobnym magazynem

**StudioAgentConflict** — Rdzen nie rozstrzyga, kto ma racje; niesie prawde o tym, co sie stalo

**AodMute** — Wyciszenie jest bytem rdzenia, nie stanem jednego okna: Operator wyciszajacy nakladke w jednej powloce ma miec cisze w kazdej nastepnej. Zniesienie idzie ta sama komenda aod.mute.set z polem muted rownym falszowi — jednym ruchem, bez pytania o potwierdzenie

**AodMute.ScopeName** — Cisza, przy ktorej Operator nie wie, CO milczy, jest gorsza od braku wyciszenia

**AodMute.EndsAt** — Puste znaczy wyciszenie trwajace do zniesienia reka Operatora — tak stoi wyciszenie kontekstowe i wyciszenie klasy zdarzen

**AodMute.DeviceId** — Wyciszenie obowiazuje wszystkie powloki Operatora niezaleznie od tego pola; pole mowi, skad przyszlo

**AodSignal** — Nosnik dla klas, ktorych rdzen nie widzi wlasna telemetria: wyniku kontroli jakosci, harmonogramu przebiegow automatyk i powtarzalnosci czynnosci Operatora (rozdz. 3.2 i 3.3 opracowania). Bez niego trzy z szesciu klas nie maja czego wyciszac ani z czego zrobic sugestii

**AodSignal.OccurrenceCount** — Prog powtarzalnosci czynnosci recznej z rozdz. 3.4 liczy sie z tego pola; puste znaczy jedno wystapienie

**MemoryDisable** — Wylacza sie albo pojedynczy wpis (entryId), albo caly poziom pamieci (level) — w zasiegu: calkiem, w srodowisku, w projekcie, w module, w parze modulow albo w karcie sesji. WYLACZENIE NIE KASUJE TRESCI i to jest jego roznica wobec memory.delete: wpis wylaczony nie wchodzi do kontekstu, ale zostaje na miejscu i wraca w calosci zniesieniem wylaczenia. Sterowanie stoi po stronie Operatora w konfiguracji — zakres ustawien „Pamiec"

**MemoryDisable.Level** — Puste znaczy wylaczenie pojedynczego wpisu wskazanego polem entryId. Zadanie bez obu pol nie ma czego wylaczyc i konczy sie odmowa

**MemoryDisabledEntry** — Wpis znika z kontekstu, nie z pamieci — tresc stoi nietknieta i wraca zniesieniem wylaczenia

**NotificationAction** — Wykonania NIE prowadzi rodzina notification.*: pozycja wskazuje komende rodziny wlasciwej i jej ladunek

**ConnectionHelloRequest.Token** — Puste znaczy polaczenie niezwiazane — bramka wtedy nie wie, kto po drugiej stronie

**ConnectionHelloResponse.LoginRequired** — Puste znaczy, ze rdzen nastawy nie zna. To nie jest bramka: pole niczego nie odmawia, mowi klientowi, czy okno logowania ma sie pokazac

**HomeEnterResponse.Presence** — Puste, gdy zadna sesja nie trwa

**SpeechAvailabilityGetRequest** — Brak silnika to ODPOWIEDZ, nie awaria — klient nie pokazuje mikrofonu, ktory nic nie nagra

**SpeechAvailabilityGetResponse** — Brak silnika to ODPOWIEDZ, nie awaria — klient nie pokazuje mikrofonu, ktory nic nie nagra

**SpeechAvailabilityGetResponse.SynthesisAvailable** — Osobne pole od available, bo odsluch jedzie innym lancuchem niz dyktowanie: piper z glosem .onnx albo espeak-ng, nie python z faster-whisper

**SpeechAvailabilityGetResponse.SynthesisReason** — Puste gdy odsluch gotowy

**SpeechTranscribeRequest** — Dzwiek NIE OPUSZCZA maszyny, na ktorej stoi silnik

**SpeechTranscribeResponse** — Dzwiek NIE OPUSZCZA maszyny, na ktorej stoi silnik

**SpeechTranscribeResponse.Transcript** — Pusty JEST dozwolony, gdy processed=true — cisza i szum to prawidlowy wynik pomiaru

**WorkspaceEnterRequest.WindowId** — Okno rozmowy nie znika przy zmianie modulu, tylko rekonfiguruje kontekst

**SessionFocusRequest** — Asystent przestawia ognisko OPERATOROWI polem targetClientId — inaczej jego posuniecia nie byly widoczne na ekranie

**SessionFocusRequest.TargetClientId** — Uzywa tego ASYSTENT dzialajacy za Operatora: bez tego pola przestawialby ognisko sobie, a ekran Operatora stalby w miejscu — czyli praca dzialaby sie "gdzies w tle"

**SessionFocusResponse** — Asystent przestawia ognisko OPERATOROWI polem targetClientId — inaczej jego posuniecia nie byly widoczne na ekranie

**WorkspaceAgentAssignRequest.Scope** — Zasieg szerszy czyni przypisanie obowiazujacym takze w innych projektach

**ResearchReportBuildResponse.FromModel** — Dzis klient nie ma jak odroznic streszczenia od komunikatu bledu, bo rdzen zbiera z kanalu same fragmenty tekstu i zapisuje je jako sekcje

**TranslateTargetAddRequest.ChannelId** — Operator nie mial dotad jak wskazac, ktory model tlumaczy

**DesignAssetGenerateRequest** — Tresc trafia do magazynu rdzenia pod suma kontrolna, wiec zasob nie zalezy od zadnego pliku zewnetrznego. Brak kanalu obrazowego albo brak poswiadczenia to ODMOWA NAZYWAJACA BRAK — rdzen nigdy nie zaklada zasobu bez bajtow obrazu

**DesignAssetGenerateRequest.ChannelId** — Brak bierze pierwszy czynny kanal obrazowy konta; brak takiego kanalu to odmowa nazywajaca brak, nigdy obraz zastepczy

**DesignAssetGenerateResponse** — Tresc trafia do magazynu rdzenia pod suma kontrolna, wiec zasob nie zalezy od zadnego pliku zewnetrznego. Brak kanalu obrazowego albo brak poswiadczenia to ODMOWA NAZYWAJACA BRAK — rdzen nigdy nie zaklada zasobu bez bajtow obrazu

**AgentCreateRequest.MemoryLevels** — Pominiete znaczy "bez zmiany" (przy zakladaniu: cztery poziomy). LISTA PUSTA `[]` JEST ZADANIEM WYLACZENIA PAMIECI i rozni sie od pominiecia pola — dlatego wylaczenie nie potrzebuje wlasnej wartosci wyliczenia

**AgentUpdateRequest.MemoryLevels** — Pominiete znaczy "bez zmiany" (przy zakladaniu: cztery poziomy). LISTA PUSTA `[]` JEST ZADANIEM WYLACZENIA PAMIECI i rozni sie od pominiecia pola — dlatego wylaczenie nie potrzebuje wlasnej wartosci wyliczenia

**AgentListRequest.ProjectId** — Podany zawezaja wykaz do ekspertow widocznych w tym projekcie: wszystkich `global` oraz tych `project`, ktore sa do niego przypisane. Pominiety znaczy wykaz biblioteki w calosci — okno Agent Buildera musi widziec takze ekspertow projektowych, inaczej nie daloby sie ich poprawic

**AgentLayerSetRequest** — Warstwa DOPISUJE sie do promptu systemowego, nie zastepuje go — trybu podania sie tu nie wybiera

**AgentLayerSetResponse** — Warstwa DOPISUJE sie do promptu systemowego, nie zastepuje go — trybu podania sie tu nie wybiera

**TranslateBacktranslationRunRequest.ChannelId** — Operator nie mial dotad jak wskazac, ktory model tlumaczy

**MemoryListResponse.DisabledEntries** — Do kontekstu nie wchodza i nie stoja w polu entries; wykaz jest zawsze tablica, bo brak wylaczen to fakt, nie cisza

**MemoryDetachRequest** — Odpiecie NIE JEST wylaczeniem: zweza zasieg samego wpisu, a nie wstrzymuje go w cudzym projekcie ani module — do tego sluzy memory.disable.set

**MemoryDetachResponse** — Odpiecie NIE JEST wylaczeniem: zweza zasieg samego wpisu, a nie wstrzymuje go w cudzym projekcie ani module — do tego sluzy memory.disable.set

**ComponentAssignRequest** — UWAGA: znaczenie przypisania nie zostalo rozstrzygniete przez ustalenie budowy

**ComponentAssignResponse** — UWAGA: znaczenie przypisania nie zostalo rozstrzygniete przez ustalenie budowy

**ExtensionInstallRequest** — UWAGA: znaczenie instalacji — pobranie paczki, zarejestrowanie adresu czy zapis punktu dostepu — nie zostalo rozstrzygniete przez ustalenie budowy

**ExtensionInstallRequest.Origin** — Rozstrzyga stan wyjsciowy: `danaco` staje wlaczone, `personal` wylaczone (rozdz. 5.2)

**ExtensionInstallResponse** — UWAGA: znaczenie instalacji — pobranie paczki, zarejestrowanie adresu czy zapis punktu dostepu — nie zostalo rozstrzygniete przez ustalenie budowy

**ExtensionUninstallRequest** — UWAGA: znaczenie odinstalowania jest zwiazane z nierozstrzygnietym znaczeniem instalacji

**ExtensionUninstallResponse** — UWAGA: znaczenie odinstalowania jest zwiazane z nierozstrzygnietym znaczeniem instalacji

**AdvisorConsultRequest** — Rada NIE JEST WIAZACA i nie wykonuje pracy za pytajacego — odpowiedzialnosc za wynik zostaje przy agencie, ktory pyta. Kim jest doradca, rozstrzyga RDZEN z danych okna, nie zadanie: pole requestedAdvisor jest PROSBA i podlega sufitowi sily. Prosba o model silniejszy od kanalu pytajacego bez wczesniejszego wskazania Operatora jest ODMAWIANA, a nie zamieniana po cichu na slabszego doradce — model z wlasnej inicjatywy po model silniejszy nie siega (rozstrzygniecie ustalenie budowy z 14.08.2026). Kazda odbyta konsultacja rozglasza sie zdarzeniem advisor.consulted

**AdvisorConsultRequest.WindowId** — Z NIEGO rdzen bierze kanal pytajacego i wskazanie Operatora — zadanie tych dwoch rzeczy nie podaje i podac nie moze

**AdvisorConsultResponse** — Rada NIE JEST WIAZACA i nie wykonuje pracy za pytajacego — odpowiedzialnosc za wynik zostaje przy agencie, ktory pyta. Kim jest doradca, rozstrzyga RDZEN z danych okna, nie zadanie: pole requestedAdvisor jest PROSBA i podlega sufitowi sily. Prosba o model silniejszy od kanalu pytajacego bez wczesniejszego wskazania Operatora jest ODMAWIANA, a nie zamieniana po cichu na slabszego doradce — model z wlasnej inicjatywy po model silniejszy nie siega (rozstrzygniecie ustalenie budowy z 14.08.2026). Kazda odbyta konsultacja rozglasza sie zdarzeniem advisor.consulted

**AdvisorConsultResponse.AdvisorChannel** — Rada nigdy nie wraca sama — bez wskazania doradcy dalaby sie podac za odpowiedz wlasna agenta

**ConfigExplainGetRequest** — Transparentnosc, nie bramka

**ConfigExplainGetResponse** — Transparentnosc, nie bramka

**ConfigWindowOpenRequest** — UWAGA: wykaz okien obiecuje te komende, nie mowiac, co rdzen ma przy niej robic poza podaniem zakresu

**ConfigWindowOpenResponse** — UWAGA: wykaz okien obiecuje te komende, nie mowiac, co rdzen ma przy niej robic poza podaniem zakresu

**TerminalOutputReadRequest** — Komenda `terminal.command.exec` konczy sie w chwili STARTU procesu (kompilacja trwa dluzej niz kazde sensowne oczekiwanie na odpowiedz, a rozlaczenie klienta nie ma prawa jej przerwac), wiec wyniku niesc nie moze i nigdy nie bedzie mogla. Tedy model czyta, co polecenie wypisalo. Zrodlem jest ten sam dziennik zbiorczego wyjscia, na ktorym stoi `terminal.output.stream` — drugiej pompy nie ma. Historia zyje jeden bieg rdzenia: po restarcie wyjscie jest puste i odpowiedz mowi to wprost

**TerminalOutputReadRequest.WaitMs** — Gorna granica to 60000 — dluzsze czekanie trzymaloby zadanie gniazda na czas, ktorego zadne polaczenie nie przezyje. Czekanie NIE jest warunkiem odpowiedzi: proces wciaz biegnacy oddaje wyjscie dotychczasowe i stan running

**TerminalOutputReadResponse** — Komenda `terminal.command.exec` konczy sie w chwili STARTU procesu (kompilacja trwa dluzej niz kazde sensowne oczekiwanie na odpowiedz, a rozlaczenie klienta nie ma prawa jej przerwac), wiec wyniku niesc nie moze i nigdy nie bedzie mogla. Tedy model czyta, co polecenie wypisalo. Zrodlem jest ten sam dziennik zbiorczego wyjscia, na ktorym stoi `terminal.output.stream` — drugiej pompy nie ma. Historia zyje jeden bieg rdzenia: po restarcie wyjscie jest puste i odpowiedz mowi to wprost

**TerminalOutputReadResponse.Truncated** — Prawda znaczy, ze poczatek zostal odciety — ostatnie bajty sa zawsze cale

**ModelChannelSetRequest** — Kanaly zaklada i zmienia rodzina channel.*; ta komenda wybiera jeden z zalozonych

**ModelChannelSetResponse** — Kanaly zaklada i zmienia rodzina channel.*; ta komenda wybiera jeden z zalozonych

**AodSuggestionRequest** — UWAGA: wykaz okien obiecuje te pozycje slowem podpowiedz, nie rozstrzygajac, czy jest to odczyt Operatora, czy zdarzenie wypychane przez rdzen

**AodSuggestionResponse** — UWAGA: wykaz okien obiecuje te pozycje slowem podpowiedz, nie rozstrzygajac, czy jest to odczyt Operatora, czy zdarzenie wypychane przez rdzen

**AuthRegisterRequest** — Konto powstaje w stanie NIEPOTWIERDZONYM i pozostaje w nim do chwili potwierdzenia adresu komenda auth.verify — dopiero potwierdzenie wydaje urzadzeniu token dostepu. Sesji ta komenda NIE zaklada. Wykonalna tylko raz; potem odmawia trwale. Powtorzenie hasla jest sprawa formularza klienta, nie kontraktu

**AuthRegisterRequest.Email** — Weryfikowany przy rejestracji, pelni pozniej funkcje metody logowania i JEDYNEJ drogi odzyskania konta — drugiego adresu do tych celow platforma nie prowadzi

**AuthRegisterResponse** — Konto powstaje w stanie NIEPOTWIERDZONYM i pozostaje w nim do chwili potwierdzenia adresu komenda auth.verify — dopiero potwierdzenie wydaje urzadzeniu token dostepu. Sesji ta komenda NIE zaklada. Wykonalna tylko raz; potem odmawia trwale. Powtorzenie hasla jest sprawa formularza klienta, nie kontraktu

**AuthRegisterResponse.PendingVerification** — Prawda oznacza, ze list z droga potwierdzenia zostal wyslany na podany adres

**AuthVerifyRequest** — Przenosi konto ze stanu niepotwierdzonego do potwierdzonego, po czym wydaje urzadzeniu token dostepu — to jest moment, w ktorym Operator wchodzi do platformy po raz pierwszy. Droga potwierdzenia przychodzi listem na adres podany przy rejestracji i wygasa; wygasla droga odmawia i pozwala poprosic o nowa

**AuthVerifyRequest.KeepSignedIn** — Nie jest bramka: znosi powtarzanie logowania, niczego nie blokuje

**AuthVerifyResponse** — Przenosi konto ze stanu niepotwierdzonego do potwierdzonego, po czym wydaje urzadzeniu token dostepu — to jest moment, w ktorym Operator wchodzi do platformy po raz pierwszy. Droga potwierdzenia przychodzi listem na adres podany przy rejestracji i wygasa; wygasla droga odmawia i pozwala poprosic o nowa

**AuthRecoverRequest** — Serwer wysyla na wskazany adres droge potwierdzenia tozsamosci; nowe haslo ustawia sie komenda auth.reset. Odpowiedz nie zdradza, czy adres pasuje do konta — inaczej komenda mowilaby obcemu, jaki adres ma Operator

**AuthRecoverResponse** — Serwer wysyla na wskazany adres droge potwierdzenia tozsamosci; nowe haslo ustawia sie komenda auth.reset. Odpowiedz nie zdradza, czy adres pasuje do konta — inaczej komenda mowilaby obcemu, jaki adres ma Operator

**AuthRecoverResponse.Sent** — Prawda niezaleznie od tego, czy adres pasuje do konta — wartosc nie jest odpowiedzia na pytanie o istnienie konta

**AuthResetRequest** — Zastepuje dotychczasowy skrot hasla przechowywany poza baza danych i uniewaznia tokeny dostepu wydane przed odzyskaniem — kazde powiazane urzadzenie loguje sie ponownie, a zmiana idzie w swiat zdarzeniem device.changed. Nowej encji Konto nie tworzy: zmienia sie wylacznie material uwierzytelniajacy

**AuthResetResponse** — Zastepuje dotychczasowy skrot hasla przechowywany poza baza danych i uniewaznia tokeny dostepu wydane przed odzyskaniem — kazde powiazane urzadzenie loguje sie ponownie, a zmiana idzie w swiat zdarzeniem device.changed. Nowej encji Konto nie tworzy: zmienia sie wylacznie material uwierzytelniajacy

**AuthLoginRequest** — Metodami dzialajacymi dzis sa HASLO i PIN; hello czeka na pochodzenie z domena po https. Login niesie sie przy metodzie password, bo konto ma nazwe ustalenie budowy nadana przy rejestracji; metody wlasciwe urzadzeniu tozsamosci nie potrzebuja, bo wskazuje ja material na urzadzeniu. Odmowa wraca kodem bledu not_authenticated. Proba nieudana NIE odmawia nastepnej — nakłada na nia ZWLOKE rosnaca wykladniczo (250 ms, 500, 1000, 2000, 4000, dalej rowno 5000 ms), zerowana przy pierwszym udanym wejsciu. Operator, ktory pomylil haslo trzy razy, wchodzi za czwartym; tylko czeka. Zadnego progu prob i zadnej odmowy 'za duzo prob' tu nie ma i nie bedzie — to byloby bramkowanie (rozstrzygniecie ustalenie budowy z 14.08.2026)

**AuthLoginRequest.Method** — Dzialajace dzis: password i pin. Metoda hello wymaga, by interfejs byl podany z pochodzenia z domena po https — pod http://127.0.0.1 jest niedostepna, bo WebAuthn wywodzi rp_id z pochodzenia dokumentu i zadna wartosc nie da sie podstawic z JavaScriptu

**AuthLoginRequest.Login** — Wymagana dla metody password; metody wlasciwe urzadzeniu (pin, hello) jej nie potrzebuja

**AuthLoginRequest.KeepSignedIn** — Nie jest bramka: znosi powtarzanie logowania, niczego nie blokuje

**AuthLoginResponse** — Metodami dzialajacymi dzis sa HASLO i PIN; hello czeka na pochodzenie z domena po https. Login niesie sie przy metodzie password, bo konto ma nazwe ustalenie budowy nadana przy rejestracji; metody wlasciwe urzadzeniu tozsamosci nie potrzebuja, bo wskazuje ja material na urzadzeniu. Odmowa wraca kodem bledu not_authenticated. Proba nieudana NIE odmawia nastepnej — nakłada na nia ZWLOKE rosnaca wykladniczo (250 ms, 500, 1000, 2000, 4000, dalej rowno 5000 ms), zerowana przy pierwszym udanym wejsciu. Operator, ktory pomylil haslo trzy razy, wchodzi za czwartym; tylko czeka. Zadnego progu prob i zadnej odmowy 'za duzo prob' tu nie ma i nie bedzie — to byloby bramkowanie (rozstrzygniecie ustalenie budowy z 14.08.2026)

**AuthMethodAddRequest** — Czynnosc USTAWIEN, wykonalna po zalogowaniu; hasla ta komenda nie zaklada, bo kotwica powstaje przy auth.register. Klucz Hello nie opuszcza TPM urzadzenia, wiec na kazdej maszynie zaklada sie go osobno — i tylko wtedy, gdy interfejs jest podawany z pochodzenia z domena po https. Dzis zalozyc mozna PIN

**AuthMethodAddRequest.Kind** — Zalozenie klucza hello wymaga pochodzenia z domena po https — pod http://127.0.0.1 rejestracja WebAuthn nie ma poprawnego rp_id

**AuthMethodAddResponse** — Czynnosc USTAWIEN, wykonalna po zalogowaniu; hasla ta komenda nie zaklada, bo kotwica powstaje przy auth.register. Klucz Hello nie opuszcza TPM urzadzenia, wiec na kazdej maszynie zaklada sie go osobno — i tylko wtedy, gdy interfejs jest podawany z pochodzenia z domena po https. Dzis zalozyc mozna PIN

**AuthMethodRemoveRequest** — Ani ostatniej metody, ani hasla zdjac sie nie da — haslo jest kotwica bramki

**AuthMethodRemoveResponse** — Ani ostatniej metody, ani hasla zdjac sie nie da — haslo jest kotwica bramki

**AuthTokenRefreshRequest** — W Danaco HUB dzieje sie to samo przy kazdym zadaniu HTTP; tu jest komenda, bo transportem jest jedno gniazdo WebSocket

**AuthTokenRefreshResponse** — W Danaco HUB dzieje sie to samo przy kazdym zadaniu HTTP; tu jest komenda, bo transportem jest jedno gniazdo WebSocket

**DeviceListRequest** — Konto jest jedno, urzadzen dowolnie wiele: komputery, telefony, tablety. Kazde niesie wlasny token dostepu, wiec wykaz jest miejscem, w ktorym Operator widzi, co ma dostep do platformy

**DeviceListResponse** — Konto jest jedno, urzadzen dowolnie wiele: komputery, telefony, tablety. Kazde niesie wlasny token dostepu, wiec wykaz jest miejscem, w ktorym Operator widzi, co ma dostep do platformy

**DeviceRevokeRequest** — Zmiana idzie do wszystkich polaczonych urzadzen zdarzeniem device.changed, wiec Operator widzi skutek natychmiast na pozostalych ekranach. Tej samej drogi uzywa odzyskanie konta, ktore uniewaznia tokeny wydane przed zmiana hasla

**DeviceRevokeResponse** — Zmiana idzie do wszystkich polaczonych urzadzen zdarzeniem device.changed, wiec Operator widzi skutek natychmiast na pozostalych ekranach. Tej samej drogi uzywa odzyskanie konta, ktore uniewaznia tokeny wydane przed zmiana hasla

**AppsDeploymentListRequest** — Bez tej komendy historia wdrozen znikala po odswiezeniu okna, bo klient odbudowywal ja wylacznie ze zdarzen biezacej sesji

**AppsDeploymentListResponse** — Bez tej komendy historia wdrozen znikala po odswiezeniu okna, bo klient odbudowywal ja wylacznie ze zdarzen biezacej sesji

**AppsArchitectureGetRequest** — Odpowiednik odczytu dla apps.architecture.define — bez niego okno po odswiezeniu nie wie, co Operator wczesniej zdefiniowal

**AppsArchitectureGetResponse** — Odpowiednik odczytu dla apps.architecture.define — bez niego okno po odswiezeniu nie wie, co Operator wczesniej zdefiniowal

**AppsWorkspaceListRequest** — Odpowiednik odczytu dla apps.workspace.update

**AppsWorkspaceListResponse** — Odpowiednik odczytu dla apps.workspace.update

**AppsWorkspaceListResponse.Files** — Kształt DeveloperFile — ten sam, którym apps.workspace.update oddaje plik po zapisie; druga struktura byłaby drugą prawdą o tym samym bycie

**DesignAssetUploadRequest** — Tresc wchodzi do magazynu rdzenia pod suma kontrolna — zasob przestaje zalezec od pliku, ktory Operator moze nadpisac albo skasowac. Uzyj, gdy zasob JUZ ISTNIEJE: Operator ma plik albo wskazuje odsylacz. Zasob majacy dopiero powstac zamawia sie komenda design.asset.generate

**DesignAssetUploadResponse** — Tresc wchodzi do magazynu rdzenia pod suma kontrolna — zasob przestaje zalezec od pliku, ktory Operator moze nadpisac albo skasowac. Uzyj, gdy zasob JUZ ISTNIEJE: Operator ma plik albo wskazuje odsylacz. Zasob majacy dopiero powstac zamawia sie komenda design.asset.generate

**DesignAssetRemoveRequest** — Rodzaj zmiany deleted zdarzenia design.asset.changed nie mial dotad zadnego nadawcy

**DesignAssetRemoveResponse** — Rodzaj zmiany deleted zdarzenia design.asset.changed nie mial dotad zadnego nadawcy

**ImageTransformRequest** — Wynikiem jest NOWY zasob w magazynie — zrodlo zostaje nietkniete, wiec model moze probowac bez ryzyka

**ImageTransformRequest.WindowId** — Model zna je z wlasnego zasiegu i podaje, zeby zasob byl WIDOCZNY dla Operatora w Assets Panel. Brak nie wstrzymuje czynnosci: bajty i tak trafiaja do magazynu pod suma kontrolna, ale zasob nie pojawi sie w wykazie okna

**ImageTransformResponse** — Wynikiem jest NOWY zasob w magazynie — zrodlo zostaje nietkniete, wiec model moze probowac bez ryzyka

**ImageAdjustRequest** — To jest RETUSZ, o ktory model prosi w rozmowie

**ImageAdjustRequest.WindowId** — Model zna je z wlasnego zasiegu i podaje, zeby zasob byl WIDOCZNY dla Operatora w Assets Panel. Brak nie wstrzymuje czynnosci: bajty i tak trafiaja do magazynu pod suma kontrolna, ale zasob nie pojawi sie w wykazie okna

**ImageAdjustResponse** — To jest RETUSZ, o ktory model prosi w rozmowie

**ImageConvertRequest** — Sluzy przygotowaniu zasobu do wydania — model uzywa jej przed osadzeniem grafiki na stronie

**ImageConvertRequest.WindowId** — Model zna je z wlasnego zasiegu i podaje, zeby zasob byl WIDOCZNY dla Operatora w Assets Panel. Brak nie wstrzymuje czynnosci: bajty i tak trafiaja do magazynu pod suma kontrolna, ale zasob nie pojawi sie w wykazie okna

**ImageConvertResponse** — Sluzy przygotowaniu zasobu do wydania — model uzywa jej przed osadzeniem grafiki na stronie

**MediaTranscodeRequest** — Wynikiem jest nowy zasob w magazynie

**MediaTranscodeRequest.WindowId** — Model zna je z wlasnego zasiegu i podaje, zeby zasob byl WIDOCZNY dla Operatora w Assets Panel. Brak nie wstrzymuje czynnosci: bajty i tak trafiaja do magazynu pod suma kontrolna, ale zasob nie pojawi sie w wykazie okna

**MediaTranscodeResponse** — Wynikiem jest nowy zasob w magazynie

**DocumentConvertRequest** — Model uzywa jej, gdy Operator prosi o dokument, a nie o tekst w oknie

**DocumentConvertRequest.WindowId** — Model zna je z wlasnego zasiegu i podaje, zeby zasob byl WIDOCZNY dla Operatora w Assets Panel. Brak nie wstrzymuje czynnosci: bajty i tak trafiaja do magazynu pod suma kontrolna, ale zasob nie pojawi sie w wykazie okna

**DocumentConvertResponse** — Model uzywa jej, gdy Operator prosi o dokument, a nie o tekst w oknie

**DocumentTextExtractRequest** — Model uzywa jej, zeby PRZECZYTAC to, co dostal jako plik

**DocumentTextExtractResponse** — Model uzywa jej, zeby PRZECZYTAC to, co dostal jako plik

**ArchivePackRequest** — Sluzy wydaniu pracy Operatorowi jednym plikiem

**ArchivePackRequest.WindowId** — Model zna je z wlasnego zasiegu i podaje, zeby zasob byl WIDOCZNY dla Operatora w Assets Panel. Brak nie wstrzymuje czynnosci: bajty i tak trafiaja do magazynu pod suma kontrolna, ale zasob nie pojawi sie w wykazie okna

**ArchivePackResponse** — Sluzy wydaniu pracy Operatorowi jednym plikiem

**ArchiveUnpackRequest.WindowId** — Model zna je z wlasnego zasiegu i podaje, zeby zasob byl WIDOCZNY dla Operatora w Assets Panel. Brak nie wstrzymuje czynnosci: bajty i tak trafiaja do magazynu pod suma kontrolna, ale zasob nie pojawi sie w wykazie okna

**MailAccountListRequest** — Bez skonfigurowanej skrzynki rodzina mail.* odmawia, nazywajac brak — rdzen poczty nie zmysla

**MailAccountListResponse** — Bez skonfigurowanej skrzynki rodzina mail.* odmawia, nazywajac brak — rdzen poczty nie zmysla

**MailMessageListRequest** — Uzyj do odnalezienia sprawy, o ktorej mowi Operator —na przykład listu od konkretnej osoby z dzisiaj

**MailMessageListResponse** — Uzyj do odnalezienia sprawy, o ktorej mowi Operator —na przykład listu od konkretnej osoby z dzisiaj

**MailMessageGetRequest** — Uzyj, gdy masz przeanalizowac list, a nie tylko go odnalezc. Zalaczniki trafiaja do magazynu rdzenia i wracaja jako zasoby, wiec model moze je dalej przetworzyc narzedziami obrazu i dokumentow

**MailMessageGetResponse** — Uzyj, gdy masz przeanalizowac list, a nie tylko go odnalezc. Zalaczniki trafiaja do magazynu rdzenia i wracaja jako zasoby, wiec model moze je dalej przetworzyc narzedziami obrazu i dokumentow

**MailDraftSaveRequest** — Szkic jest krokiem POSREDNIM: Operator widzi go w swojej poczcie, zanim cokolwiek wyjdzie w swiat

**MailDraftSaveResponse** — Szkic jest krokiem POSREDNIM: Operator widzi go w swojej poczcie, zanim cokolwiek wyjdzie w swiat

**MailSendRequest** — TO JEST JEDYNA KOMENDA RODZINY, KTOREJ SKUTEK WYCHODZI POZA MASZYNE OPERATORA i nie da sie go cofnac

**MailSendResponse** — TO JEST JEDYNA KOMENDA RODZINY, KTOREJ SKUTEK WYCHODZI POZA MASZYNE OPERATORA i nie da sie go cofnac

**MailAccountAddRequest** — Rdzen nie zaklada konta pocztowego i nie stawia serwera — bierze skrzynke, ktora juz istnieje: na urzadzeniu albo w chmurze. Poswiadczenie idzie do sejfu, nie do bazy

**MailAccountAddResponse** — Rdzen nie zaklada konta pocztowego i nie stawia serwera — bierze skrzynke, ktora juz istnieje: na urzadzeniu albo w chmurze. Poswiadczenie idzie do sejfu, nie do bazy

**MailAccountDiscoverRequest** — Niczego nie podpina i nie siega po haslo — oddaje to, co znalazl, zeby Operator nie przepisywal nastaw recznie

**MailAccountDiscoverResponse** — Niczego nie podpina i nie siega po haslo — oddaje to, co znalazl, zeby Operator nie przepisywal nastaw recznie

**MailAccountRemoveRequest** — Skrzynki u dostawcy nie tyka

**MailAccountRemoveResponse** — Skrzynki u dostawcy nie tyka

**ImageUpscaleRequest** — Zastepuje wyspecjalizowane modele powiekszajace; bez zainstalowanego silnika odmawia, nazywajac brak — nigdy nie oddaje zwyklego rozciagniecia jako powiekszenia

**ImageUpscaleResponse** — Zastepuje wyspecjalizowane modele powiekszajace; bez zainstalowanego silnika odmawia, nazywajac brak — nigdy nie oddaje zwyklego rozciagniecia jako powiekszenia

**ImageBackgroundRemoveRequest** — Zastepuje wyspecjalizowane narzedzia wycinania; bez silnika odmawia, nazywajac brak

**ImageBackgroundRemoveResponse** — Zastepuje wyspecjalizowane narzedzia wycinania; bez silnika odmawia, nazywajac brak

**KnowledgeIndexRequest** — Wyszukiwanie po slowach juz dziala (library.file.search); to jest droga do wyszukiwania po SENSIE

**KnowledgeIndexResponse** — Wyszukiwanie po slowach juz dziala (library.file.search); to jest droga do wyszukiwania po SENSIE

**KnowledgeSearchRequest.Rerank** — Podobienstwo wektorow jest pierwszym przebiegiem: pytanie i fragment licza sie osobno i spotykaja dopiero jako dwie liczby. Przesiew jest przebiegiem drugim — czyta pytanie RAZEM z fragmentem i uklada kandydatow na nowo. Kosztuje wczytanie drugiego modelu, wiec wchodzi na zadanie, a nie zawsze

**KnowledgeSearchRequest.RerankCandidates** — Liczba wieksza od `limit` jest tu sensem rzeczy: przesiew moze wyniesc na czolo fragment, ktory po samych wektorach byl dwudziesty. Bez `rerank` nie znaczy nic

**KnowledgeSearchResponse.Reranked** — Bez tego pola nie da sie odroznic odpowiedzi przesianej od odpowiedzi z samego pierwszego przebiegu, a `score` w obu ma ten sam ksztalt

**KnowledgeImageSearchRequest** — Rodzina knowledge.* prowadzi dotad wylacznie tekst; ta komenda jest jej osia obrazu i przeglada obrazy biblioteki Operatora

**KnowledgeImageSearchResponse** — Rodzina knowledge.* prowadzi dotad wylacznie tekst; ta komenda jest jej osia obrazu i przeglada obrazy biblioteki Operatora

**KnowledgeImageSearchResponse.Examined** — Wynik pusty przy zerze znaczy brak obrazow, a nie brak trafienia — dwie rozne odpowiedzi, ktorych bez tej liczby nie da sie rozroznic

**SessionToolAttachRequest** — Definicji eksperta NIE RUSZA (rozstrzygniecie ustalenie budowy, rozdz. 5): dolozenie zyje w stanie sesji, przezywa rozlaczenie klienta i konczy sie wraz z sesja albo z usunieciem rozmowy

**SessionToolAttachResponse** — Definicji eksperta NIE RUSZA (rozstrzygniecie ustalenie budowy, rozdz. 5): dolozenie zyje w stanie sesji, przezywa rozlaczenie klienta i konczy sie wraz z sesja albo z usunieciem rozmowy

**SessionToolDetachRequest** — Zestaw wraca do podstawy z definicji eksperta; sama definicja nie zmienia sie ani o joto, bo nigdy nie byla zmieniana

**SessionToolDetachResponse** — Zestaw wraca do podstawy z definicji eksperta; sama definicja nie zmienia sie ani o joto, bo nigdy nie byla zmieniana

**SessionToolListRequest** — Zestaw narzedzi tury to definicja eksperta PLUS te dolozenia — bez tego odczytu druga polowa zestawu bylaby niewidoczna

**SessionToolListResponse** — Zestaw narzedzi tury to definicja eksperta PLUS te dolozenia — bez tego odczytu druga polowa zestawu bylaby niewidoczna

**ToolsCatalogListRequest** — Wykaz liczy setki pozycji, wiec zadanie niesie zawezenie tekstem, rodzajem i grupa; dopasowanie idzie takze SRODKIEM nazwy, bo przy przedrostkach zrodla szukanie od poczatku byloby bezuzyteczne (rozdz. 7.2)

**ToolsCatalogListResponse** — Wykaz liczy setki pozycji, wiec zadanie niesie zawezenie tekstem, rodzajem i grupa; dopasowanie idzie takze SRODKIEM nazwy, bo przy przedrostkach zrodla szukanie od poczatku byloby bezuzyteczne (rozdz. 7.2)

**AgentConnectorListRequest** — Odpowiednik `agent.plugin.list` dla drugiego z dwoch bytow

**AgentConnectorListResponse** — Odpowiednik `agent.plugin.list` dla drugiego z dwoch bytow

**AgentVersionGetRequest** — `AgentVersion` niesie etykiete, opis zmiany, sprawce i sume kontrolna, ale NIE niesie tresci — bez tej komendy podgladu wersji nie ma z czego zlozyc, a przywrocenie jest jedynym sposobem zobaczenia, co w wersji stalo

**AgentVersionGetResponse** — `AgentVersion` niesie etykiete, opis zmiany, sprawce i sume kontrolna, ale NIE niesie tresci — bez tej komendy podgladu wersji nie ma z czego zlozyc, a przywrocenie jest jedynym sposobem zobaczenia, co w wersji stalo

**AgentAssignmentListRequest** — `workspace.agent.assign` zapisuje przynaleznosc, ale zadna komenda nie odczytuje jej OD STRONY EKSPERTA, wiec biblioteka nie ma skad wziac liczby

**AgentAssignmentListResponse** — `workspace.agent.assign` zapisuje przynaleznosc, ale zadna komenda nie odczytuje jej OD STRONY EKSPERTA, wiec biblioteka nie ma skad wziac liczby

**AgentModulesSetRequest** — LISTA PUSTA `[]` ZNACZY BRAK OGRANICZENIA, czyli dostepnosc wszedzie, i jest to stan wyjsciowy; pominiecie pola znaczy „bez zmiany". Ta sama zasada co przy `memoryLevels`, gdzie zbior pusty tez niesie znaczenie wlasne

**AgentModulesSetResponse** — LISTA PUSTA `[]` ZNACZY BRAK OGRANICZENIA, czyli dostepnosc wszedzie, i jest to stan wyjsciowy; pominiecie pola znaczy „bez zmiany". Ta sama zasada co przy `memoryLevels`, gdzie zbior pusty tez niesie znaczenie wlasne

**AgentIsolationSetRequest** — Wartosc obowiazuje wszedzie, gdzie ekspert dziala, dopoki nie nadpisze jej regula zapisana na poziomie zasiegu bardziej szczegolowym

**AgentIsolationSetResponse** — Wartosc obowiazuje wszedzie, gdzie ekspert dziala, dopoki nie nadpisze jej regula zapisana na poziomie zasiegu bardziej szczegolowym

**AgentPolicyGetRequest** — Transparentnosc, nie bramka — komenda niczego nie zapisuje i niczego nie rozstrzyga

**AgentPolicyGetResponse** — Transparentnosc, nie bramka — komenda niczego nie zapisuje i niczego nie rozstrzyga

**AgentSubagentSetRequest** — Granica pietnastu jest parametrem technicznym platformy, nie decyzja produktowa — ta sama, ktora zna `subagent.spawn`

**AgentSubagentSetResponse** — Granica pietnastu jest parametrem technicznym platformy, nie decyzja produktowa — ta sama, ktora zna `subagent.spawn`

**AgentPermissionRemoveRequest** — `agent.permission.set` wylacznie USTAWIA wartosc wpisu, wiec po pierwszym zawezeniu wiersz zakresu zostaje w wykazie na zawsze, choc jego wartosc wrocila do przyznanej

**AgentPermissionRemoveResponse** — `agent.permission.set` wylacznie USTAWIA wartosc wpisu, wiec po pierwszym zawezeniu wiersz zakresu zostaje w wykazie na zawsze, choc jego wartosc wrocila do przyznanej

**ChannelCheckRequest** — Narzedzie pomocnicze, nie bramka: wynik nie warunkuje zapisu eksperta ani wyslania tury

**ChannelCheckResponse** — Narzedzie pomocnicze, nie bramka: wynik nie warunkuje zapisu eksperta ani wyslania tury

**ChannelCredentialStatusRequest** — NIE ODDAJE TRESCI POSWIADCZENIA I ODDAWAC JEJ NIE MOZE

**ChannelCredentialStatusResponse** — NIE ODDAJE TRESCI POSWIADCZENIA I ODDAWAC JEJ NIE MOZE

**ExtensionSearchRequest** — Dopelnia extension.list, ktory zawezal wylacznie rodzajem i stanem zainstalowania; dopasowanie idzie takze SRODKIEM nazwy, bo pozycje niosa przedrostki zrodla

**ExtensionSearchResponse** — Dopelnia extension.list, ktory zawezal wylacznie rodzajem i stanem zainstalowania; dopasowanie idzie takze SRODKIEM nazwy, bo pozycje niosa przedrostki zrodla

**ExtensionDetailGetRequest** — Extension niesie sam naglowek pozycji, wiec karta szczegolow nie ma dzis z czego powstac

**ExtensionDetailGetResponse** — Extension niesie sam naglowek pozycji, wiec karta szczegolow nie ma dzis z czego powstac

**ExtensionCollectionApplyRequest** — Wynik jest bilansem, nie potwierdzeniem: pozycja odrzucona wraca z powodem, zamiast znikac po cichu

**ExtensionCollectionApplyResponse** — Wynik jest bilansem, nie potwierdzeniem: pozycja odrzucona wraca z powodem, zamiast znikac po cichu

**ExtensionRegistryListRequest** — Wyliczenie ExtensionOrigin ma dzis dwie wartosci, wiec rejestr organizacji nie ma jak sie w katalogu pokazac

**ExtensionRegistryListResponse** — Wyliczenie ExtensionOrigin ma dzis dwie wartosci, wiec rejestr organizacji nie ma jak sie w katalogu pokazac

**ExtensionUpdateCheckRequest** — Extension niesie wlasna wersje, ale nie wersje dostepna, wiec wskaznika aktualizacji nie ma dzis z czego zlozyc

**ExtensionUpdateCheckResponse** — Extension niesie wlasna wersje, ale nie wersje dostepna, wiec wskaznika aktualizacji nie ma dzis z czego zlozyc

**ExtensionVersionPinRequest** — Wersja jest dzis polem opisowym, wiec nie ma czego przypiac

**ExtensionVersionPinResponse** — Wersja jest dzis polem opisowym, wiec nie ma czego przypiac

**ExtensionVersionRollbackRequest** — Wykonuje sie od razu; ustawienie extension.rollback.confirm wlacza potwierdzenie, ktore nie warunkuje wykonania

**ExtensionVersionRollbackResponse** — Wykonuje sie od razu; ustawienie extension.rollback.confirm wlacza potwierdzenie, ktore nie warunkuje wykonania

**ExtensionBundleInstallRequest** — Wynik jest bilansem: pozycja odrzucona wraca z powodem

**ExtensionBundleInstallResponse** — Wynik jest bilansem: pozycja odrzucona wraca z powodem

**ExtensionHistoryListRequest** — Rejestr niesie stan biezacy, nie droge, ktora do niego doprowadzila

**ExtensionHistoryListResponse** — Rejestr niesie stan biezacy, nie droge, ktora do niego doprowadzila

**ExtensionAdminBulkRequest** — Wynik jest bilansem przyjetych i odrzuconych, nie pojedynczym potwierdzeniem

**ExtensionAdminBulkResponse** — Wynik jest bilansem przyjetych i odrzuconych, nie pojedynczym potwierdzeniem

**ExtensionToolListRequest** — Odpowiedz na tools/list, resources/list i prompts/list nie ma dzis w kontrakcie ani komendy, ani struktury, w ktora mialaby wejsc

**ExtensionToolListResponse** — Odpowiedz na tools/list, resources/list i prompts/list nie ma dzis w kontrakcie ani komendy, ani struktury, w ktora mialaby wejsc

**ExtensionToolCallRequest** — Wywolanie jest probne: idzie poza kontekstem eksperta i niczego mu nie przypisuje

**ExtensionToolCallResponse** — Wywolanie jest probne: idzie poza kontekstem eksperta i niczego mu nie przypisuje

**ExtensionTransportSetRequest** — Transport da sie dzis wpisac wylacznie do nieprzezroczystej konfiguracji, wiec rdzen i okno moga rozumiec to pole inaczej

**ExtensionTransportSetResponse** — Transport da sie dzis wpisac wylacznie do nieprzezroczystej konfiguracji, wiec rdzen i okno moga rozumiec to pole inaczej

**ExtensionCredentialBindRequest** — Moduł operuje na referencjach; tresc poswiadczenia nigdy nie opuszcza rdzenia i nie wchodzi w to zadanie

**ExtensionCredentialBindResponse** — Moduł operuje na referencjach; tresc poswiadczenia nigdy nie opuszcza rdzenia i nie wchodzi w to zadanie

**ExtensionSecretListRequest** — Oddaje wylacznie odwolania i metryke; tresc sekretu nie opuszcza rdzenia

**ExtensionSecretListResponse** — Oddaje wylacznie odwolania i metryke; tresc sekretu nie opuszcza rdzenia

**ExtensionPermissionListRequest** — Extension niesie nieprzezroczysta konfiguracje, ktora deklaracja uprawnien nie jest

**ExtensionPermissionListResponse** — Extension niesie nieprzezroczysta konfiguracje, ktora deklaracja uprawnien nie jest

**ExtensionPermissionGrantRequest** — Nadanie jest jedyna kontrola; brak nadania nie wstrzymuje instalacji ani wlaczenia (zasada zero blokad)

**ExtensionPermissionGrantResponse** — Nadanie jest jedyna kontrola; brak nadania nie wstrzymuje instalacji ani wlaczenia (zasada zero blokad)

**ExtensionManifestScanRequest** — Wynik jest SYGNALEM, nie brama: nie wstrzymuje instalacji ani wlaczenia

**ExtensionManifestScanResponse** — Wynik jest SYGNALEM, nie brama: nie wstrzymuje instalacji ani wlaczenia

**AppsStageListRequest** — Odpowiednik odczytu dla zdarzenia apps.build.changed — bez niego wykaz etapow zaczyna pusty przy kazdym wejsciu w modul, choc rdzen trzyma etapy w bazie

**AppsStageListResponse** — Odpowiednik odczytu dla zdarzenia apps.build.changed — bez niego wykaz etapow zaczyna pusty przy kazdym wejsciu w modul, choc rdzen trzyma etapy w bazie

**AppsStageSaveRequest** — AppStage.ownerAgentId jedzie dzis wylacznie w kierunku od rdzenia, wiec przypisania wykonawcy nie ma jak zlecic

**AppsStageSaveResponse** — AppStage.ownerAgentId jedzie dzis wylacznie w kierunku od rdzenia, wiec przypisania wykonawcy nie ma jak zlecic

**AppsProductLinkListRequest** — Powiazania nie sa czynne domyslnie — kazde jest swiadoma decyzja Operatora

**AppsProductLinkListResponse** — Powiazania nie sa czynne domyslnie — kazde jest swiadoma decyzja Operatora

**AppsArchitectureValidateRequest** — Dzis pole validationIssues jest zapisywane i oddawane bez wyliczenia, wiec jego pustka nie orzeka o poprawnosci ukladu. Zastrzezenia sa OSTRZEZENIEM, nie brama — nie wstrzymuja pracy w warsztatach

**AppsArchitectureValidateResponse** — Dzis pole validationIssues jest zapisywane i oddawane bez wyliczenia, wiec jego pustka nie orzeka o poprawnosci ukladu. Zastrzezenia sa OSTRZEZENIEM, nie brama — nie wstrzymuja pracy w warsztatach

**AppsArchitectureVersionListRequest** — AppArchitecture.version jest liczba, ale historii kontrakt nie niesie, wiec porownania nie ma z czym zestawic

**AppsArchitectureVersionListResponse** — AppArchitecture.version jest liczba, ale historii kontrakt nie niesie, wiec porownania nie ma z czym zestawic

**AppsPreviewStartRequest** — Podglad na zywo jest narzedziem warstwy 1 Frontend Workspace, a rdzen nie stawia dzis serwera produktu

**AppsPreviewStartResponse** — Podglad na zywo jest narzedziem warstwy 1 Frontend Workspace, a rdzen nie stawia dzis serwera produktu

**AppsEndpointListRequest** — AppComponent niesie apiContract jako JEDEN tekst, wiec punkt koncowy nie jest dzis bytem, ktory dalby sie wyliczyc

**AppsEndpointListResponse** — AppComponent niesie apiContract jako JEDEN tekst, wiec punkt koncowy nie jest dzis bytem, ktory dalby sie wyliczyc

**AppsEnvironmentListRequest** — Wyliczenie AppDeployEnvironment ma trzy wartosci stale, a opracowanie zada dowolnej liczby srodowisk konfigurowanych przez Operatora

**AppsEnvironmentListResponse** — Wyliczenie AppDeployEnvironment ma trzy wartosci stale, a opracowanie zada dowolnej liczby srodowisk konfigurowanych przez Operatora

**AppsEnvironmentVariableListRequest** — Wartosc sekretna oddawana jest jako referencja, nigdy jako tresc

**AppsEnvironmentVariableListResponse** — Wartosc sekretna oddawana jest jako referencja, nigdy jako tresc

**AppsServiceLogReadRequest** — Podglad na zywo idzie zdarzeniem apps.service.log

**AppsServiceLogReadResponse** — Podglad na zywo idzie zdarzeniem apps.service.log

**AppsArtifactListRequest** — AppDeployment niesie wersje i odnosnik logu, ale artefakt nie jest dzis bytem umowy

**AppsArtifactListResponse** — AppDeployment niesie wersje i odnosnik logu, ale artefakt nie jest dzis bytem umowy

**AppsDeploymentLogReadRequest** — AppDeployment niesie logRef — sam odnosnik, bez tresci — wiec podgladu logu nie ma dzis czym wypelnic. Podglad na zywo idzie zdarzeniem apps.deployment.log

**AppsDeploymentLogReadResponse** — AppDeployment niesie logRef — sam odnosnik, bez tresci — wiec podgladu logu nie ma dzis czym wypelnic. Podglad na zywo idzie zdarzeniem apps.deployment.log

**AppsPackageValidateRequest** — Zastrzezenia sa OSTRZEZENIEM, nie brama: nie wstrzymuja publikacji

**AppsPackageValidateResponse** — Zastrzezenia sa OSTRZEZENIEM, nie brama: nie wstrzymuja publikacji

**AppsPackageSignRequest** — Klucz wydawcy lezy w warstwie sekretow i nie wchodzi w to zadanie ani w odpowiedz

**AppsPackageSignResponse** — Klucz wydawcy lezy w warstwie sekretow i nie wchodzi w to zadanie ani w odpowiedz

**SpeechAudioUploadRequest** — Jest to brak platformowy, nie brak jednego okna: speech.transcribe bierze audioRef, czyli SCIEZKE PLIKU na maszynie silnika, a nagranie z mikrofonu istnieje wylacznie jako bajty w pamieci karty. Ten sam brak wstrzymuje dyktowanie w oknie komunikacji (client/src/okno-komunikacji/dyktowanie/dostarczenie-nagrania.ts) i mikrofon Voice Console — jedna komenda otwiera obie drogi naraz. Dzwiek nie opuszcza maszyny rdzenia: bajty ida do magazynu nagran, nie do sieci

**SpeechAudioUploadRequest.Retain** — Falsz kasuje je po rozpoznaniu; wartosc domyslna bierze ustawienie mowa_zapis_nagran

**SpeechAudioUploadResponse** — Jest to brak platformowy, nie brak jednego okna: speech.transcribe bierze audioRef, czyli SCIEZKE PLIKU na maszynie silnika, a nagranie z mikrofonu istnieje wylacznie jako bajty w pamieci karty. Ten sam brak wstrzymuje dyktowanie w oknie komunikacji (client/src/okno-komunikacji/dyktowanie/dostarczenie-nagrania.ts) i mikrofon Voice Console — jedna komenda otwiera obie drogi naraz. Dzwiek nie opuszcza maszyny rdzenia: bajty ida do magazynu nagran, nie do sieci

**SpeechAudioFetchRequest** — Domyka pare do speech.audio.upload i do translate.speech.synthesize: dzis rdzen oddaje odnosniki (AssistantActivityEntry.audioRef, AssistantVoiceCommandResponse.speechRef), lecz nie ma komendy, ktora by po nie siegnela — odslucha nie ma czym wykonac

**SpeechAudioFetchResponse** — Domyka pare do speech.audio.upload i do translate.speech.synthesize: dzis rdzen oddaje odnosniki (AssistantActivityEntry.audioRef, AssistantVoiceCommandResponse.speechRef), lecz nie ma komendy, ktora by po nie siegnela — odslucha nie ma czym wykonac

**SpeechWakeGetRequest** — Brak modelu frazy jest ODPOWIEDZIA, nie awaria — tak samo jak brak silnika w speech.availability.get

**SpeechWakeGetResponse** — Brak modelu frazy jest ODPOWIEDZIA, nie awaria — tak samo jak brak silnika w speech.availability.get

**SpeechWakeSetRequest** — Pola pominiete zostaja bez zmian

**SpeechWakeSetResponse** — Pola pominiete zostaja bez zmian

**SpeechListenStartRequest** — Rdzen slucha strumienia i oglasza zdarzeniami: czesciowa transkrypcje oraz wykrycie frazy wybudzajacej. Nasluch nie jest bramka — Operator zatrzymuje go speech.listen.stop, a rdzen nie zatrzymuje go sam

**SpeechListenStartResponse** — Rdzen slucha strumienia i oglasza zdarzeniami: czesciowa transkrypcje oraz wykrycie frazy wybudzajacej. Nasluch nie jest bramka — Operator zatrzymuje go speech.listen.stop, a rdzen nie zatrzymuje go sam

**SpeechListenStopRequest** — Zatrzymanie nasluchu, ktorego nie ma, nie jest bledem

**SpeechListenStopResponse** — Zatrzymanie nasluchu, ktorego nie ma, nie jest bledem

**AssistantActivityFlagRequest** — Dzis wyroznienie wpisu nie mialoby gdzie zamieszkac: AssistantActivityEntry nie niesie takiego pola, wiec okno pokazywaloby wyroznienie, ktorego rdzen nie pamieta

**AssistantActivityFlagResponse** — Dzis wyroznienie wpisu nie mialoby gdzie zamieszkac: AssistantActivityEntry nie niesie takiego pola, wiec okno pokazywaloby wyroznienie, ktorego rdzen nie pamieta

**MemoryRetentionGetRequest** — Odczyt do pary z memory.retention.set

**MemoryRetentionGetResponse** — Odczyt do pary z memory.retention.set

**MemoryRetentionSetRequest** — Zasada obejmuje zapisy kolejne, a wpisow zastanych nie rusza wstecz — inaczej zmiana reguly kasowalaby ustalenia, na ktore Operator sie nie umawial

**MemoryRetentionSetResponse** — Zasada obejmuje zapisy kolejne, a wpisow zastanych nie rusza wstecz — inaczej zmiana reguly kasowalaby ustalenia, na ktore Operator sie nie umawial

**MemoryContextListRequest** — Dzis memory.toggle przestawia poziomy zasiegu i to jedyny kontekst przelaczalny; nazwy zestawu nie ma gdzie odlozyc

**MemoryContextListResponse** — Dzis memory.toggle przestawia poziomy zasiegu i to jedyny kontekst przelaczalny; nazwy zestawu nie ma gdzie odlozyc

**MemoryContextSaveRequest** — Puste contextId zaklada nowy, podane zmienia istniejacy — tak samo jak memory.set

**MemoryContextSaveResponse** — Puste contextId zaklada nowy, podane zmienia istniejacy — tak samo jak memory.set

**MemoryContextActivateRequest** — Aktywacja nie kasuje kontekstu poprzedniego — wraca sie do niego tym samym wywolaniem

**MemoryContextActivateResponse** — Aktywacja nie kasuje kontekstu poprzedniego — wraca sie do niego tym samym wywolaniem

**MemoryContextDeleteRequest** — Wpisow pamieci nie kasuje — kontekst jest zestawem wskazan, a nie ustalenie budowy tresci

**MemoryContextDeleteResponse** — Wpisow pamieci nie kasuje — kontekst jest zestawem wskazan, a nie ustalenie budowy tresci

**ToolsScopeListRequest** — Dzis takiego zakresu nie ma dla profilu: agent.connector.* i agent.permission.set dotycza eksperta modulu Agents i nie siegaja profilu asystenta

**ToolsScopeListResponse** — Dzis takiego zakresu nie ma dla profilu: agent.connector.* i agent.permission.set dotycza eksperta modulu Agents i nie siegaja profilu asystenta

**ToolsScopeSetRequest** — Zakres jest nastawa zasiegu, nie bramka: stanem wyjsciowym jest pelny dostep bez limitu, zgodnie z zasada zero blokad

**ToolsScopeSetResponse** — Zakres jest nastawa zasiegu, nie bramka: stanem wyjsciowym jest pelny dostep bez limitu, zgodnie z zasada zero blokad

**ClipboardListRequest** — Schowek jest bytem klienta, ktorego rdzen dzis nie zna, wiec historia ginie wraz z karta — komenda daje jej trwalosc

**ClipboardListResponse** — Schowek jest bytem klienta, ktorego rdzen dzis nie zna, wiec historia ginie wraz z karta — komenda daje jej trwalosc

**ClipboardPushRequest** — Powtorzenie tresci identycznej nie mnozy wpisow — podnosi wpis zastany na czolo wykazu

**ClipboardPushResponse** — Powtorzenie tresci identycznej nie mnozy wpisow — podnosi wpis zastany na czolo wykazu

**ClipboardPinRequest** — Przypiety nie wygasa wraz z zasada retencji historii

**ClipboardPinResponse** — Przypiety nie wygasa wraz z zasada retencji historii

**SnippetListRequest** — Skrot dziala we wszystkich polach tekstowych platformy, wiec slownik nalezy do rdzenia, nie do jednego okna

**SnippetListResponse** — Skrot dziala we wszystkich polach tekstowych platformy, wiec slownik nalezy do rdzenia, nie do jednego okna

**SnippetSetRequest** — Puste snippetId zaklada nowy, podane zmienia istniejacy

**SnippetSetResponse** — Puste snippetId zaklada nowy, podane zmienia istniejacy

**LauncherHotkeyGetRequest** — Skrot globalny rejestruje powloka programu okiennego, wiec brak wsparcia powloki jest ODPOWIEDZIA, nie awaria — klient nie obiecuje wtedy skrotu, ktory nikogo nie obudzi

**LauncherHotkeyGetResponse** — Skrot globalny rejestruje powloka programu okiennego, wiec brak wsparcia powloki jest ODPOWIEDZIA, nie awaria — klient nie obiecuje wtedy skrotu, ktory nikogo nie obudzi

**LauncherHotkeyGetResponse.Hotkey** — Ctrl+Shift+Space; pusty znaczy brak nastawy

**LauncherHotkeySetRequest** — Skrot zajety przez inny program nie jest bledem zapisu — zapis zostaje, a odpowiedz mowi, ze rejestracja sie nie udala

**LauncherHotkeySetResponse** — Skrot zajety przez inny program nie jest bledem zapisu — zapis zostaje, a odpowiedz mowi, ze rejestracja sie nie udala

**ContextUsageGetRequest** — ZALEZNOSC: wymaga tokenizatora jako podsystemu rdzenia — bez niego liczba tokenow bylaby wartoscia wzieta znikad. Do czasu jego zbudowania komenda ma zwracac available falsz wraz z powodem, tak jak speech.availability.get przy braku silnika

**ContextUsageGetResponse** — ZALEZNOSC: wymaga tokenizatora jako podsystemu rdzenia — bez niego liczba tokenow bylaby wartoscia wzieta znikad. Do czasu jego zbudowania komenda ma zwracac available falsz wraz z powodem, tak jak speech.availability.get przy braku silnika

**DesignAssetContentGetRequest** — Pole uri zasobu jest sciezka w systemie plikow rdzenia (magazyn oddaje filepath), wiec przegladarka nie wczyta spod niego niczego — take wtedy, gdy zasob powstal bez zarzutu. Komenda dotyczy KAZDEGO zasobu magazynu, nie tylko obrazu: jeden magazyn obsluguje rodziny design.*, document.*, media.* i archive.*. Rdzen odmawia zasobu, ktorego nie zna, i zasobu, ktorego tresci nie ma pod suma kontrolna — nigdy nie oddaje bajtow zastepczych

**DesignAssetContentGetResponse** — Pole uri zasobu jest sciezka w systemie plikow rdzenia (magazyn oddaje filepath), wiec przegladarka nie wczyta spod niego niczego — take wtedy, gdy zasob powstal bez zarzutu. Komenda dotyczy KAZDEGO zasobu magazynu, nie tylko obrazu: jeden magazyn obsluguje rodziny design.*, document.*, media.* i archive.*. Rdzen odmawia zasobu, ktorego nie zna, i zasobu, ktorego tresci nie ma pod suma kontrolna — nigdy nie oddaje bajtow zastepczych

**DesignAssetExportRequest** — Rozni sie od image.convert celem: konwersja zaklada NOWY zasob w magazynie i tam sie konczy, eksport oddaje bajty gotowe do zapisania poza produktem

**DesignAssetExportResponse** — Rozni sie od image.convert celem: konwersja zaklada NOWY zasob w magazynie i tam sie konczy, eksport oddaje bajty gotowe do zapisania poza produktem

**DesignAssetExportBatchRequest** — Pokrywa eksport zbiorczy Assets Panel oraz zestawy rozmiarow kampanii; odmowa jednego zasobu NIE wstrzymuje pozostalych — wynik niesie bilans przyjetych i odrzuconych

**DesignAssetExportBatchResponse** — Pokrywa eksport zbiorczy Assets Panel oraz zestawy rozmiarow kampanii; odmowa jednego zasobu NIE wstrzymuje pozostalych — wynik niesie bilans przyjetych i odrzuconych

**DesignCollectionCreateRequest** — Etykiety zasobu juz sa (design.asset.tag.set), ale kolekcja jest bytem osobnym: ma nazwe, opis i porzadek, a etykieta jest tylko slowem

**DesignCollectionCreateResponse** — Etykiety zasobu juz sa (design.asset.tag.set), ale kolekcja jest bytem osobnym: ma nazwe, opis i porzadek, a etykieta jest tylko slowem

**DesignCollectionAssignRequest** — Zestaw jest DOKLADKA albo ODJECIEM, nie zastapieniem — inaczej niz przy etykietach, bo kolekcja bywa duza i przepisywanie jej w calosci przy kazdej zmianie jest droga do zgubienia zawartosci

**DesignCollectionAssignResponse** — Zestaw jest DOKLADKA albo ODJECIEM, nie zastapieniem — inaczej niz przy etykietach, bo kolekcja bywa duza i przepisywanie jej w calosci przy kazdej zmianie jest droga do zgubienia zawartosci

**DesignCollectionListRequest** — Bez tego kolekcja zalozona nie mialaby jak wrocic na ekran po odswiezeniu

**DesignCollectionListResponse** — Bez tego kolekcja zalozona nie mialaby jak wrocic na ekran po odswiezeniu

**DesignPromptTemplateSaveRequest** — Dzis historia promptow i szablony zyja w oknie do zamkniecia karty przegladarki i tyle o nich wiadomo

**DesignPromptTemplateSaveResponse** — Dzis historia promptow i szablony zyja w oknie do zamkniecia karty przegladarki i tyle o nich wiadomo

**DesignPromptHistoryListRequest** — Domyka tez prowenancje: dzis zasob niesie pole promptId, ktorego rdzen NIE wypelnia, bo nie ma przekladu klucza wiersza promptu na kod kontraktu — ta komenda ten przeklad wnosi

**DesignPromptHistoryListResponse** — Domyka tez prowenancje: dzis zasob niesie pole promptId, ktorego rdzen NIE wypelnia, bo nie ma przekladu klucza wiersza promptu na kod kontraktu — ta komenda ten przeklad wnosi

**DesignBoardVersionRestoreRequest** — Przywrocenie ZAKLADA nowa wersje z ukladu sprzed przywrocenia, zeby cofniecie sie samo nie kasowalo stanu, ktory Operator wlasnie porzucil

**DesignBoardVersionRestoreResponse** — Przywrocenie ZAKLADA nowa wersje z ukladu sprzed przywrocenia, zeby cofniecie sie samo nie kasowalo stanu, ktory Operator wlasnie porzucil

**DesignBoardExportRequest** — Dzis kompozycja jezdzi do rdzenia i z powrotem jako uklad warstw i nie ma drogi wyjscia poza rdzen

**DesignBoardExportResponse** — Dzis kompozycja jezdzi do rdzenia i z powrotem jako uklad warstw i nie ma drogi wyjscia poza rdzen

**DesignAnnotationSetRequest** — Pole note warstwy niesie JEDNO zdanie bez autora i bez watku; opracowanie wymaga watkow i oznaczen osob, a tego jedno pole nie unosi

**DesignAnnotationSetResponse** — Pole note warstwy niesie JEDNO zdanie bez autora i bez watku; opracowanie wymaga watkow i oznaczen osob, a tego jedno pole nie unosi

**DesignPresenceReportRequest** — Rdzen rozglasza je pozostalym zdarzeniem design.board.presence. Zgloszenie jest ULOTNE — nie zapisuje sie w bazie, bo polozenie kursora sprzed godziny nie jest wiedza o niczym

**DesignPresenceReportResponse** — Rdzen rozglasza je pozostalym zdarzeniem design.board.presence. Zgloszenie jest ULOTNE — nie zapisuje sie w bazie, bo polozenie kursora sprzed godziny nie jest wiedza o niczym

**DesignTokensetSaveRequest** — Dzis Tokens & System Panel czyta zetony z motywu obowiazujacego i nie ma ich gdzie odlozyc — kontrakt nie zna bytu zestawu zetonow

**DesignTokensetSaveResponse** — Dzis Tokens & System Panel czyta zetony z motywu obowiazujacego i nie ma ich gdzie odlozyc — kontrakt nie zna bytu zestawu zetonow

**DesignTokensetExportRequest** — Klient sklada dzis zmienne CSS, SCSS, konfiguracje Tailwind i modul JavaScript sam i nie potrzebuje do tego rdzenia; ta komenda jest potrzebna dla postaci, ktorych przegladarka zlozyc nie moze, oraz dla wydania idacego DO INNEGO MODULU zamiast do pliku

**DesignTokensetExportResponse** — Klient sklada dzis zmienne CSS, SCSS, konfiguracje Tailwind i modul JavaScript sam i nie potrzebuje do tego rdzenia; ta komenda jest potrzebna dla postaci, ktorych przegladarka zlozyc nie moze, oraz dla wydania idacego DO INNEGO MODULU zamiast do pliku

**DesignTokensetImportRequest** — Rdzen NIE nadpisuje motywu produktu — motyw jest wlasnoscia powloki; import zaklada byt obok niego i oddaje roznice wobec zetonow wskazanego motywu

**DesignTokensetImportResponse** — Rdzen NIE nadpisuje motywu produktu — motyw jest wlasnoscia powloki; import zaklada byt obok niego i oddaje roznice wobec zetonow wskazanego motywu

**DesignStyleguidePublishRequest** — Przegladarka sklada przewodnik sama i oddaje go plikiem, ale wydania go do Library albo Studio nie ma czym zlecic

**DesignStyleguidePublishResponse** — Przegladarka sklada przewodnik sama i oddaje go plikiem, ale wydania go do Library albo Studio nie ma czym zlecic

**ImageVectorizeRequest** — Rodzina image.* zna przeksztalcenie geometryczne, poprawke, konwersje formatu, powiekszenie i wyciecie tla — zamiany rastra na sciezki nie zna zadna z nich, a jest to funkcja opracowania modulu

**ImageVectorizeResponse** — Rodzina image.* zna przeksztalcenie geometryczne, poprawke, konwersje formatu, powiekszenie i wyciecie tla — zamiany rastra na sciezki nie zna zadna z nich, a jest to funkcja opracowania modulu

**ImageLayersSplitRequest** — Zasila rozdzielenie generacji na warstwy oraz zaznaczanie obiektu; bez silnika segmentacji ODMAWIA, nazywajac brak — nigdy nie oddaje calego obrazu jako jednej warstwy udajac rozklad

**ImageLayersSplitResponse** — Zasila rozdzielenie generacji na warstwy oraz zaznaczanie obiektu; bez silnika segmentacji ODMAWIA, nazywajac brak — nigdy nie oddaje calego obrazu jako jednej warstwy udajac rozklad

**ImageComposeRequest** — Pokrywa znak wodny, branding wsadowy, osadzenie w ramce urzadzenia i warstwy rastrowe — cztery funkcje opracowania, ktore wszystkie potrzebuja jednego: zlozenia obrazu na obrazie

**ImageComposeResponse** — Pokrywa znak wodny, branding wsadowy, osadzenie w ramce urzadzenia i warstwy rastrowe — cztery funkcje opracowania, ktore wszystkie potrzebuja jednego: zlozenia obrazu na obrazie

**RoundtableDebateGetRequest** — Obszar nie ma dzis zadnej komendy odczytu, wiec okno otwarte w trakcie debaty zna wylacznie to, co uslyszalo zdarzeniem od swojego otwarcia

**RoundtableDebateGetResponse** — Obszar nie ma dzis zadnej komendy odczytu, wiec okno otwarte w trakcie debaty zna wylacznie to, co uslyszalo zdarzeniem od swojego otwarcia

**StudioDocumentFormatSetRequest** — Dzis format nadaje rdzen przy wczytaniu i zadna komenda go nie zmienia: studio.document.save przyjmuje documentId, content, title i createVersion, ale nie format

**StudioDocumentFormatSetResponse** — Dzis format nadaje rdzen przy wczytaniu i zadna komenda go nie zmienia: studio.document.save przyjmuje documentId, content, title i createVersion, ale nie format

**StudioCommentAddRequest.Author** — Model zaklada komentarz podpisany jako model i ta droga, nie druga

**StudioCommentResolveRequest** — Watek nie znika: rozwiazanie jest stanem, nie usunieciem

**StudioCommentResolveResponse** — Watek nie znika: rozwiazanie jest stanem, nie usunieciem

**StudioTrackingSetRequest** — Przy wlaczonym sledzeniu kazdy zapis odklada wstawienia i usuniecia jako zmiany do decyzji, zamiast nadpisywac tresc

**StudioTrackingSetResponse** — Przy wlaczonym sledzeniu kazdy zapis odklada wstawienia i usuniecia jako zmiany do decyzji, zamiast nadpisywac tresc

**StudioTrackingDecideRequest** — Decyzja zapisuje sie w tresci dokumentu i zaklada wersje

**StudioTrackingDecideResponse** — Decyzja zapisuje sie w tresci dokumentu i zaklada wersje

**StudioAssetEmbedRequest** — Realizuje powiazanie Design do Studio

**StudioAssetEmbedResponse** — Realizuje powiazanie Design do Studio

**StudioOperationDeleteRequest** — Operacji fabrycznej nie usuwa — na nia odpowiada odmowa nazywajaca ten fakt

**StudioOperationDeleteResponse** — Operacji fabrycznej nie usuwa — na nia odpowiada odmowa nazywajaca ten fakt

**StudioChainRunRequest** — Przebieg prowadzi petla wykonawcza okna, a kazdy krok odklada wlasna propozycje zmiany

**StudioChainRunResponse** — Przebieg prowadzi petla wykonawcza okna, a kazdy krok odklada wlasna propozycje zmiany

**StudioBatchRunRequest** — Kazdy dokument dostaje wlasne zadanie petli; odmowa jednego nie wstrzymuje pozostalych

**StudioBatchRunResponse** — Kazdy dokument dostaje wlasne zadanie petli; odmowa jednego nie wstrzymuje pozostalych

**StudioAnnotationAddRequest** — Dzis czynnosc ta idzie droga generyczna window.action i wraca odmowa not_found, bo katalog akcji nie ma jej wiersza

**StudioAnnotationAddResponse** — Dzis czynnosc ta idzie droga generyczna window.action i wraca odmowa not_found, bo katalog akcji nie ma jej wiersza

**StudioDiffReportExportRequest** — Wynik jest zasobem magazynu rdzenia

**StudioDiffReportExportResponse** — Wynik jest zasobem magazynu rdzenia

**StudioDiffVisualRequest** — Sluzy tam, gdzie roznica tekstowa nie widzi zmiany ukladu

**StudioDiffVisualResponse** — Sluzy tam, gdzie roznica tekstowa nie widzi zmiany ukladu

**StudioSearchSemanticRequest** — Uzupelnia wyszukiwanie wzorca w studio.diff.compare

**StudioSearchSemanticResponse** — Uzupelnia wyszukiwanie wzorca w studio.diff.compare

**StudioDiffSourceRequest** — Sluzy kontroli, czy redakcja nie odeszla od zrodla

**StudioDiffSourceResponse** — Sluzy kontroli, czy redakcja nie odeszla od zrodla

**StudioProposalDecideRequest** — Dzis decyzja zapada wylacznie w kliencie i rdzen o niej nie wie, wiec propozycja zostaje w nim nierozstrzygnieta

**StudioProposalDecideResponse** — Dzis decyzja zapada wylacznie w kliencie i rdzen o niej nie wie, wiec propozycja zostaje w nim nierozstrzygnieta

**StudioBranchCreateRequest** — Sluzy prowadzeniu dwoch redakcji obok siebie

**StudioBranchCreateResponse** — Sluzy prowadzeniu dwoch redakcji obok siebie

**StudioBranchMergeRequest** — Konflikt nierozstrzygniety wraca w wyniku zamiast byc rozstrzygniety domyslem

**StudioBranchMergeResponse** — Konflikt nierozstrzygniety wraca w wyniku zamiast byc rozstrzygniety domyslem

**StudioPreviewRenderRequest** — Zwraca strony jako zasoby, wiec podglad pokazuje UKLAD, a nie sam tekst

**StudioPreviewRenderResponse** — Zwraca strony jako zasoby, wiec podglad pokazuje UKLAD, a nie sam tekst

**StudioPdfSplitRequest** — Kazda czesc jest osobnym zasobem

**StudioPdfSplitResponse** — Kazda czesc jest osobnym zasobem

**StudioSecurityEncryptRequest** — Czynnosc jest jawnym, odwracalnym ustawieniem Operatora i niczego nie warunkuje

**StudioSecurityEncryptResponse** — Czynnosc jest jawnym, odwracalnym ustawieniem Operatora i niczego nie warunkuje

**StudioSecurityRedactRequest** — Czynnosc jest nieodwracalna dla wyniku, dlatego zrodlo zostaje nietkniete

**StudioSecurityRedactResponse** — Czynnosc jest nieodwracalna dla wyniku, dlatego zrodlo zostaje nietkniete

**StudioSecuritySensitiveDetectRequest** — Czynnosc CZYTA i niczego nie zmienia

**StudioSecuritySensitiveDetectResponse** — Czynnosc CZYTA i niczego nie zmienia

**StudioIngestQueueAddRequest** — Dzis kolejka jest wylacznie kliencka i ginie z odswiezeniem okna

**StudioIngestQueueAddResponse** — Dzis kolejka jest wylacznie kliencka i ginie z odswiezeniem okna

**StudioIngestRecognizeRequest** — Komenda document.text.extract robi to samo waskim wejsciem — jeden jezyk, bez silnika, bez progu i bez ukladu

**StudioIngestRecognizeResponse** — Komenda document.text.extract robi to samo waskim wejsciem — jeden jezyk, bez silnika, bez progu i bez ukladu

**StudioIngestUrlRequest** — Migawke strony oddaje browser.snapshot.get, ale wymaga okna modulu Browser i nie prowadzi do dokumentu Studia

**StudioIngestUrlResponse** — Migawke strony oddaje browser.snapshot.get, ale wymaga okna modulu Browser i nie prowadzi do dokumentu Studia

**StudioIngestDeviceListRequest** — Bez tego wykazu wybor urzadzenia nie ma z czego powstac

**StudioIngestDeviceListResponse** — Bez tego wykazu wybor urzadzenia nie ma z czego powstac

**StudioStyleSaveRequest** — Zmiana stylu przestawia wszystkie miejsca dokumentu, ktore go uzywaja

**StudioStyleSaveResponse** — Zmiana stylu przestawia wszystkie miejsca dokumentu, ktore go uzywaja

**StudioStyleDeleteRequest** — Stylu fabrycznego nie usuwa — odpowiada odmowa nazywajaca powod

**StudioStyleDeleteResponse** — Stylu fabrycznego nie usuwa — odpowiada odmowa nazywajaca powod

**StudioPageSetupSetRequest** — Zmiana formatu przelicza uklad i oddaje bilans tego, co sie nie zmiescilo

**StudioPageSetupSetResponse** — Zmiana formatu przelicza uklad i oddaje bilans tego, co sie nie zmiescilo

**StudioTableStructureEditRequest** — Szerokosci kolumn zostaja policzone, nie zerowe

**StudioTableStructureEditResponse** — Szerokosci kolumn zostaja policzone, nie zerowe

**StudioDocumentImportPdfRequest** — Odzyskanie jest odtworzeniem, nie odczytem, wiec odpowiedz niesie bilans; PDF ze samych skanow kieruje na rozpoznanie tekstu

**StudioDocumentImportPdfResponse** — Odzyskanie jest odtworzeniem, nie odczytem, wiec odpowiedz niesie bilans; PDF ze samych skanow kieruje na rozpoznanie tekstu

**StudioDocumentCopyRequest** — Kopia jest osobnym dokumentem, nie drugim odwolaniem do tego samego

**StudioDocumentCopyResponse** — Kopia jest osobnym dokumentem, nie drugim odwolaniem do tego samego

**StudioDocumentExportFormatRequest** — Format ubozszy niz dokument oddaje wykaz cech pominietych

**StudioDocumentExportFormatResponse** — Format ubozszy niz dokument oddaje wykaz cech pominietych

**StudioTemplateDeleteRequest** — Szablonu fabrycznego nie usuwa — odpowiada odmowa nazywajaca powod

**StudioTemplateDeleteResponse** — Szablonu fabrycznego nie usuwa — odpowiada odmowa nazywajaca powod

**StudioLockAddRequest** — Blokada obowiazuje w rdzeniu, przed dotknieciem tresci

**StudioLockAddResponse** — Blokada obowiazuje w rdzeniu, przed dotknieciem tresci

**StudioLockRemoveRequest** — Blokade zdejmuje WYLACZNIE Operator — czynnosc modelu wraca odmowa nazywajaca powod

**StudioLockRemoveResponse** — Blokade zdejmuje WYLACZNIE Operator — czynnosc modelu wraca odmowa nazywajaca powod

**StudioJournalRevertRequest** — Czynnosc bedaca podstawa pozniejszej odmawia i nazywa zaleznosc, zamiast zostawic dokument w stanie niespojnym

**StudioJournalRevertResponse** — Czynnosc bedaca podstawa pozniejszej odmawia i nazywa zaleznosc, zamiast zostawic dokument w stanie niespojnym

**StudioModelChangesRevertRequest** — Nie jest to przywrocenie wersji sprzed, bo to skasowaloby prace Operatora

**StudioModelChangesRevertResponse** — Nie jest to przywrocenie wersji sprzed, bo to skasowaloby prace Operatora

**StudioAutosaveRunRequest** — Nieudany zapis wraca nazwany, nie przemilczany

**StudioAutosaveRunResponse** — Nieudany zapis wraca nazwany, nie przemilczany

**StudioViewSetRequest** — Gdzie da sie zrobic dwojako i obie drogi maja sens, wybor nalezy do Operatora i jest jawnym, odwracalnym ustawieniem

**StudioViewSetResponse** — Gdzie da sie zrobic dwojako i obie drogi maja sens, wybor nalezy do Operatora i jest jawnym, odwracalnym ustawieniem

**StudioMarkupAddRequest** — Znakowanie modelu jest podpisane jako model

**StudioMarkupAddResponse** — Znakowanie modelu jest podpisane jako model

**StudioMarkupTypeDeleteRequest** — Rodzaju fabrycznego nie usuwa — odpowiada odmowa nazywajaca powod

**StudioMarkupTypeDeleteResponse** — Rodzaju fabrycznego nie usuwa — odpowiada odmowa nazywajaca powod

**StudioAgentsSettingsGetRequest** — Oba narzedzia sa domyslnie wylaczone

**StudioAgentsSettingsGetResponse** — Oba narzedzia sa domyslnie wylaczone

**StudioAgentsSettingsSetRequest** — Wartosci ida zasiegami rodziny config, nie osobnym magazynem

**StudioAgentsSettingsSetResponse** — Wartosci ida zasiegami rodziny config, nie osobnym magazynem

**StudioPlanCreateRequest** — Sam rozklad niczego nie uruchamia

**StudioPlanCreateResponse** — Sam rozklad niczego nie uruchamia

**StudioPlanRunRequest** — Gdy nastawa petli jest wylaczona, wraca odmowa nazywajaca brak nastawy, a nie cisza

**StudioPlanRunResponse** — Gdy nastawa petli jest wylaczona, wraca odmowa nazywajaca brak nastawy, a nie cisza

**StudioAgentsClaimRequest** — Fragment zajety przez kogo innego wraca odmowa nazywajaca wykonawce i zakres

**StudioAgentsClaimResponse** — Fragment zajety przez kogo innego wraca odmowa nazywajaca wykonawce i zakres

**StudioAgentsConflictsListRequest** — Odlozone brzmienie nie przepada

**StudioAgentsConflictsListResponse** — Odlozone brzmienie nie przepada

**TerminalSessionCloseRequest** — Dzis zamkniecie karty zyje wylacznie w widoku klienta: powloka i jej procesy biegna dalej, a rdzen o zamknieciu nie wie

**TerminalSessionCloseResponse** — Dzis zamkniecie karty zyje wylacznie w widoku klienta: powloka i jej procesy biegna dalej, a rdzen o zamknieciu nie wie

**TerminalSessionListRequest** — Rdzen odtwarza karty przy starcie, ale klient po ponownym polaczeniu nie ma jak ich zobaczyc i zaczyna wykaz od pustego

**TerminalSessionListResponse** — Rdzen odtwarza karty przy starcie, ale klient po ponownym polaczeniu nie ma jak ich zobaczyc i zaczyna wykaz od pustego

**TerminalFileReadRequest** — Bez tej komendy klient czyta manifest projektu poleceniem powloki, wiec wykrycie zadan zalezy od programu wypisujacego plik i od skladni kazdej z szesciu powlok

**TerminalFileReadResponse** — Bez tej komendy klient czyta manifest projektu poleceniem powloki, wiec wykrycie zadan zalezy od programu wypisujacego plik i od skladni kazdej z szesciu powlok

**TerminalProcessSuspendRequest** — Dzis rdzen umie proces wylacznie zakonczyc, wiec dlugie zadanie da sie tylko ubic

**TerminalProcessSuspendResponse** — Dzis rdzen umie proces wylacznie zakonczyc, wiec dlugie zadanie da sie tylko ubic

**TerminalHostSaveRequest** — Bez tej komendy ksiazka zyje jedno posiedzenie przegladarki i ginie przy odswiezeniu strony

**TerminalHostSaveResponse** — Bez tej komendy ksiazka zyje jedno posiedzenie przegladarki i ginie przy odswiezeniu strony

**TerminalHostRemoveRequest** — Karty juz otwarte do tego hosta biegna dalej

**TerminalHostRemoveResponse** — Karty juz otwarte do tego hosta biegna dalej

**TerminalScriptSaveRequest** — Bez tej komendy biblioteka zyje jedno posiedzenie, a jedyna droga jej zachowania jest wywoz do pliku

**TerminalScriptSaveResponse** — Bez tej komendy biblioteka zyje jedno posiedzenie, a jedyna droga jej zachowania jest wywoz do pliku

**TerminalScriptLintRequest** — Analize prowadza programy spoza instalki Danaco Console, wiec odpowiedz mowi wprost, czy narzedzie bylo dostepne — pusty wykaz uwag przy braku narzedzia znaczylby falszywie tresc bez zastrzezen

**TerminalScriptLintResponse** — Analize prowadza programy spoza instalki Danaco Console, wiec odpowiedz mowi wprost, czy narzedzie bylo dostepne — pusty wykaz uwag przy braku narzedzia znaczylby falszywie tresc bez zastrzezen

**TerminalTunnelOpenRequest** — Dzis tunel da sie zalozyc wylacznie poleceniem wydanym w karcie, a wtedy jego stan i przepustowosc sa niewidoczne

**TerminalTunnelOpenResponse** — Dzis tunel da sie zalozyc wylacznie poleceniem wydanym w karcie, a wtedy jego stan i przepustowosc sa niewidoczne

**TerminalKeyGenerateRequest** — Haslo klucza wchodzi odwolaniem do sejfu, nigdy trescia — zgodnie z zasada zapisana w kontrakcie przy zmiennej srodowiska

**TerminalKeyGenerateResponse** — Haslo klucza wchodzi odwolaniem do sejfu, nigdy trescia — zgodnie z zasada zapisana w kontrakcie przy zmiennej srodowiska

**TerminalKeyImportRequest** — Klucz wskazuje sie sciezka, a nie trescia: material kluczowy nie ma powodu przechodzic przez lacze

**TerminalKeyImportResponse** — Klucz wskazuje sie sciezka, a nie trescia: material kluczowy nie ma powodu przechodzic przez lacze

**TerminalKeyRemoveRequest** — Wpisy ksiazki hostow wskazujace ten klucz traca wskazanie i wracaja do klucza domyslnego konfiguracji maszyny

**TerminalKeyRemoveResponse** — Wpisy ksiazki hostow wskazujace ten klucz traca wskazanie i wracaja do klucza domyslnego konfiguracji maszyny

**TerminalWatchStartRequest** — Kontrakt daje dzis wyzwalacz plikowy automatyce, a nie karcie powloki

**TerminalWatchStartResponse** — Kontrakt daje dzis wyzwalacz plikowy automatyce, a nie karcie powloki

**TerminalWatchStopRequest** — Polecenie juz uruchomione biegnie dalej

**TerminalWatchStopResponse** — Polecenie juz uruchomione biegnie dalej

**WorkspaceAgentUnassignRequest** — Nie usuwa eksperta z biblioteki modulu Agents — znosi wylacznie jego przypisanie do tego projektu

**WorkspaceAgentUnassignResponse** — Nie usuwa eksperta z biblioteki modulu Agents — znosi wylacznie jego przypisanie do tego projektu

**WorkspaceAgentUnassignResponse.Unassigned** — Falsz znaczy, ze rdzen takiego przypisania nie mial — nie jest to blad wywolania

**WorkspaceProjectStatusSetRequest** — Droga do stanu 'paused', ktorego opracowanie wymaga, a ktorego zaden dzisiejszy uchwyt nie potrafi zapisac

**WorkspaceProjectStatusSetResponse** — Droga do stanu 'paused', ktorego opracowanie wymaga, a ktorego zaden dzisiejszy uchwyt nie potrafi zapisac

**WorkspaceInstructionsVersionRestoreRequest** — Przywrocenie zaklada wersje nowa o tresci wersji wskazanej — historia nie jest przepisywana

**WorkspaceInstructionsVersionRestoreResponse** — Przywrocenie zaklada wersje nowa o tresci wersji wskazanej — historia nie jest przepisywana

**WorkspaceTaskUpdateRequest** — Pola pominiete zostaja bez zmiany — wywolanie nie jest podmiana calego zadania

**WorkspaceTaskUpdateResponse** — Pola pominiete zostaja bez zmiany — wywolanie nie jest podmiana calego zadania

**WorkspaceTaskMoveRequest** — Osobno od workspace.task.update, bo przeciagniecie karty zmienia stan i porzadek naraz, a jest czynnoscia jednym ruchem myszy

**WorkspaceTaskMoveResponse** — Osobno od workspace.task.update, bo przeciagniecie karty zmienia stan i porzadek naraz, a jest czynnoscia jednym ruchem myszy

**WorkspaceTaskMoveResponse.WipExceeded** — Ostrzezenie, nie odmowa — platforma nie stawia blokad

**WorkspaceTaskDependencySetRequest** — Rdzen odmawia zalozenia zaleznosci domykajacej cykl, bo cykl nie da sie ulozyc w czasie

**WorkspaceTaskDependencySetResponse** — Rdzen odmawia zalozenia zaleznosci domykajacej cykl, bo cykl nie da sie ulozyc w czasie

**WorkspaceScheduleGetRequest** — Zadanie bez obu granic czasu do harmonogramu nie wchodzi

**WorkspaceScheduleGetResponse** — Zadanie bez obu granic czasu do harmonogramu nie wchodzi

**WorkspaceScheduleGetResponse.UnscheduledTaskIds** — Pole niezbedne: bez niego zadania te znikalyby z widoku bez sladu

**WorkspaceCalendarImportResponse.Skipped** — Milczace pomijanie zostawiloby Operatora z niepelnym kalendarzem bez sladu

**WorkspaceNoteSaveRequest** — Puste 'noteId' zaklada notatke nowa; podane zmienia istniejaca. Rdzen przy zapisie przelicza odnosniki tresci

**WorkspaceNoteSaveResponse** — Puste 'noteId' zaklada notatke nowa; podane zmienia istniejaca. Rdzen przy zapisie przelicza odnosniki tresci

**WorkspaceKnowledgeGraphGetRequest** — Nazwa rodziny 'workspace.knowledge' nie miesza sie z rodzina 'knowledge.*': tamta prowadzi wskaznik ZNACZENIA, ta rysuje siec ODNOSNIKOW miedzy bytami projektu

**WorkspaceKnowledgeGraphGetResponse** — Nazwa rodziny 'workspace.knowledge' nie miesza sie z rodzina 'knowledge.*': tamta prowadzi wskaznik ZNACZENIA, ta rysuje siec ODNOSNIKOW miedzy bytami projektu

**WorkspaceCanvasGetResponse.Canvas** — Pole pominiete znaczy, ze projekt nie ma jeszcze zadnej tablicy — nie jest to odmowa

**WorkspaceLibraryTextExtractRequest** — Obejmuje rozpoznanie tekstu z obrazow i skanow

**WorkspaceLibraryTextExtractResponse** — Obejmuje rozpoznanie tekstu z obrazow i skanow

**WorkspaceLibraryDuplicateListRequest** — Sam wykaz niczego nie scala

**WorkspaceLibraryDuplicateListResponse** — Sam wykaz niczego nie scala

**WorkspaceLibraryDuplicateMergeRequest** — Etykiety, kolekcje i wersje plikow scalanych przechodza na plik zachowany, zeby scalenie nie gubilo dorobku

**WorkspaceLibraryDuplicateMergeResponse** — Etykiety, kolekcje i wersje plikow scalanych przechodza na plik zachowany, zeby scalenie nie gubilo dorobku

**WorkspaceSearchProjectRequest** — Nie zastepuje library.file.search: tamta komenda przeszukuje biblioteke centralna, ta jeden projekt i wiecej niz pliki. Wyszukiwanie po ZNACZENIU prowadzi rodzina knowledge.

**WorkspaceSearchProjectResponse** — Nie zastepuje library.file.search: tamta komenda przeszukuje biblioteke centralna, ta jeden projekt i wiecej niz pliki. Wyszukiwanie po ZNACZENIU prowadzi rodzina knowledge.

**LibraryMetadataGetRequest** — Metadane osadzone czyta sie z bajtow, wiec ich odczyt jest kosztowny i wchodzi wylacznie na wyrazne zadanie

**LibraryMetadataGetResponse** — Metadane osadzone czyta sie z bajtow, wiec ich odczyt jest kosztowny i wchodzi wylacznie na wyrazne zadanie

**LibraryMetadataSetRequest** — Dzis formularz opisu nie ma dokad pojsc: library.tag.set zmienia wylacznie etykiety i kolekcje. Zapis rozglasza library.file.changed

**LibraryMetadataSetResponse** — Dzis formularz opisu nie ma dokad pojsc: library.tag.set zmienia wylacznie etykiety i kolekcje. Zapis rozglasza library.file.changed

**LibrarySchemaSetRequest** — Pole o kodzie juz istniejacym jest zmieniane, nie dublowane

**LibrarySchemaSetResponse** — Pole o kodzie juz istniejacym jest zmieniane, nie dublowane

**LibraryTagListRequest** — Bez tej komendy klient sklada slownik z etykiet plikow odczytanej strony wykazu, wiec pokazuje probke zamiast slownika

**LibraryTagListResponse** — Bez tej komendy klient sklada slownik z etykiet plikow odczytanej strony wykazu, wiec pokazuje probke zamiast slownika

**LibraryTagUpdateRequest** — Zmiana nazwy przechodzi na wszystkie zasoby noszace etykiete — inaczej powstalaby druga etykieta o tym samym znaczeniu

**LibraryTagUpdateResponse** — Zmiana nazwy przechodzi na wszystkie zasoby noszace etykiete — inaczej powstalaby druga etykieta o tym samym znaczeniu

**LibraryTagMergeRequest** — Zasoby noszace etykiete zrodlowa dostaja docelowa, a zrodlowa znika ze slownika

**LibraryTagMergeResponse** — Zasoby noszace etykiete zrodlowa dostaja docelowa, a zrodlowa znika ze slownika

**LibraryTagRemoveRequest** — Usuniecie etykiety uzywanej zada potwierdzenia, bo zdejmuje ja z zasobow, ktorych zadanie nie wymienia

**LibraryTagRemoveResponse** — Usuniecie etykiety uzywanej zada potwierdzenia, bo zdejmuje ja z zasobow, ktorych zadanie nie wymienia

**LibraryCollectionListRequest** — Bez tej komendy klient zna wylacznie identyfikatory kolekcji wyczytane z plikow, a nazwe kolekcji tylko w chwili jej zalozenia

**LibraryCollectionListResponse** — Bez tej komendy klient zna wylacznie identyfikatory kolekcji wyczytane z plikow, a nazwe kolekcji tylko w chwili jej zalozenia

**LibraryRuleSetRequest** — Zapis uruchamia przeliczenie

**LibraryRuleSetResponse** — Zapis uruchamia przeliczenie

**LibraryRuleRemoveRequest** — Zasoby przypisane przez regule kolekcji inteligentnej zostaja w kolekcji, ale przestaja byc przeliczane — inaczej usuniecie reguly oproznialoby kolekcje po cichu

**LibraryRuleRemoveResponse** — Zasoby przypisane przez regule kolekcji inteligentnej zostaja w kolekcji, ale przestaja byc przeliczane — inaczej usuniecie reguly oproznialoby kolekcje po cichu

**LibraryDuplicateScanRequest** — Rozpoznanie przyblizone zada porownania calego zbioru, wiec jest czynnoscia rdzenia, nie okna

**LibraryDuplicateScanResponse** — Rozpoznanie przyblizone zada porownania calego zbioru, wiec jest czynnoscia rdzenia, nie okna

**LibraryDuplicateMergeRequest** — Zasoby wchloniete trafiaja do archiwum, nie znikaja

**LibraryDuplicateMergeResponse** — Zasoby wchloniete trafiaja do archiwum, nie znikaja

**LibraryFixityCheckRequest** — Klient nie ma czym tego zrobic, bo nie siega po bajty zasobu

**LibraryFixityCheckResponse** — Klient nie ma czym tego zrobic, bo nie siega po bajty zasobu

**LibraryNameNormalizeRequest** — Przebieg probny pokazuje wynik bez zapisu

**LibraryNameNormalizeResponse** — Przebieg probny pokazuje wynik bez zapisu

**LibraryAuditListRequest** — Dziennik jest przyrostowy — wpisow nie da sie zmienic ani usunac ta rodzina komend

**LibraryAuditListResponse** — Dziennik jest przyrostowy — wpisow nie da sie zmienic ani usunac ta rodzina komend

**LibraryRetentionSetRequest** — Twarde usuniecie nie zachodzi samoczynnie: polityka najwyzej zglasza zasob do usuniecia

**LibraryRetentionSetResponse** — Twarde usuniecie nie zachodzi samoczynnie: polityka najwyzej zglasza zasob do usuniecia

**LibraryPreservationRunRequest.Profile** — PDF/A-2b; brak bierze odmiane z konfiguracji modulu

**LibraryFileMoveRequest** — Przeniesienie nie rusza tresci ani historii wersji

**LibraryFileMoveResponse** — Przeniesienie nie rusza tresci ani historii wersji

**LibraryFileArchiveRequest** — Czynnosc jest odwracalna komenda library.file.restore — to jest kosz repozytorium, nie usuniecie

**LibraryFileArchiveResponse** — Czynnosc jest odwracalna komenda library.file.restore — to jest kosz repozytorium, nie usuniecie

**LibraryFileDeleteRequest** — Jedyna droga utraty danych repozytorium; bez potwierdzenia jest odmowa

**LibraryFileDeleteResponse** — Jedyna droga utraty danych repozytorium; bez potwierdzenia jest odmowa

**LibraryShareCreateRequest** — Token jest jawny zgodnie z zasada jawnosci kluczy platformy

**LibraryShareCreateResponse** — Token jest jawny zgodnie z zasada jawnosci kluczy platformy

**LibraryShareRevokeRequest** — Odnosnik przestaje dzialac, wpis zostaje w dzienniku audytu

**LibraryShareRevokeResponse** — Odnosnik przestaje dzialac, wpis zostaje w dzienniku audytu

**LibraryClassifyRunRequest** — Wynik wraca sugestiami do przyjecia, nie zapisem — domyslnym zachowaniem modulu jest sugestia z akceptacja Operatora

**LibraryClassifyRunResponse** — Wynik wraca sugestiami do przyjecia, nie zapisem — domyslnym zachowaniem modulu jest sugestia z akceptacja Operatora

**LibrarySuggestionApplyRequest** — Przyjecie wykonuje czynnosc, ktora sugestia opisuje; odrzucenie zdejmuje ja z wykazu

**LibrarySuggestionApplyResponse** — Przyjecie wykonuje czynnosc, ktora sugestia opisuje; odrzucenie zdejmuje ja z wykazu

**LibraryDiffCompareRequest** — Dokumenty binarne porownuje po tekscie z nich wydobytym i mowi o tym w odpowiedzi

**LibraryDiffCompareResponse** — Dokumenty binarne porownuje po tekscie z nich wydobytym i mowi o tym w odpowiedzi

**DesignVectorPathSetRequest** — Wezly przychodza w calosci — zmiana jednego wezla idzie ta sama droga co narysowanie sciezki, zeby nie bylo dwoch prawd o jej ksztalcie

**DesignVectorPathSetResponse** — Wezly przychodza w calosci — zmiana jednego wezla idzie ta sama droga co narysowanie sciezki, zeby nie bylo dwoch prawd o jej ksztalcie

**DesignVectorShapeAddRequest** — Ksztalt powstaje od razu jako wezly, wiec da sie go dalej edytowac pioram — nie jest osobnym bytem, ktory potem trzeba zamieniac

**DesignVectorShapeAddResponse** — Ksztalt powstaje od razu jako wezly, wiec da sie go dalej edytowac pioram — nie jest osobnym bytem, ktory potem trzeba zamieniac

**DesignVectorBooleanRequest** — Sciezki zrodlowe znikaja albo zostaja wedle wskazania — bo suma dwoch ksztaltow bywa krokiem posrednim, a bywa wynikiem koncowym

**DesignVectorBooleanResponse** — Sciezki zrodlowe znikaja albo zostaja wedle wskazania — bo suma dwoch ksztaltow bywa krokiem posrednim, a bywa wynikiem koncowym

**DesignVectorTextPathRequest** — Zamiana w kontury jest nieodwracalna dla wyniku, dlatego tekst zrodlowy zostaje

**DesignVectorTextPathResponse** — Zamiana w kontury jest nieodwracalna dla wyniku, dlatego tekst zrodlowy zostaje

**DesignVectorOptimizeRequest** — Odpowiedz niesie ubytek zmierzony, zeby Operator widzial, ile naprawde ubylo, zamiast czytac obietnice

**DesignVectorOptimizeResponse** — Odpowiedz niesie ubytek zmierzony, zeby Operator widzial, ile naprawde ubylo, zamiast czytac obietnice

**DesignVectorSymbolSetRequest** — Zmiana definicji propaguje do wszystkich instancji — po to symbol jest

**DesignVectorSymbolSetResponse** — Zmiana definicji propaguje do wszystkich instancji — po to symbol jest

**DesignFrameSetRequest** — Ramka jest ekranem makiety: to ona, a nie kanwa, wyznacza obszar wydania

**DesignFrameSetResponse** — Ramka jest ekranem makiety: to ona, a nie kanwa, wyznacza obszar wydania

**DesignFrameRemoveRequest** — Warstwy ramki zostaja na kanwie — usuniecie ramki nie jest usunieciem pracy, ktora w niej lezala

**DesignFrameRemoveResponse** — Warstwy ramki zostaja na kanwie — usuniecie ramki nie jest usunieciem pracy, ktora w niej lezala

**DesignLayoutAutoRequest** — Rdzen oddaje warstwy po przeliczeniu, wiec klient nie liczy ukladu drugi raz i nie ma jak sie rozjechac

**DesignLayoutAutoResponse** — Rdzen oddaje warstwy po przeliczeniu, wiec klient nie liczy ukladu drugi raz i nie ma jak sie rozjechac

**DesignFrameResizeApplyRequest** — Osobno od design.frame.set, bo tam zmiana rozmiaru jest zapisem nastawy, a tu przeliczeniem ukladu wedle wiezi

**DesignFrameResizeApplyResponse** — Osobno od design.frame.set, bo tam zmiana rozmiaru jest zapisem nastawy, a tu przeliczeniem ukladu wedle wiezi

**DesignMockupGenerateRequest** — Wynik jest makieta do poprawienia, nie projektem koncowym, i tak sie o nim mowi

**DesignMockupGenerateResponse** — Wynik jest makieta do poprawienia, nie projektem koncowym, i tak sie o nim mowi

**DesignMockupImportRequest** — Odpowiedz mowi, czego rdzen nie rozpoznal — bilans zamiast ciszy

**DesignMockupImportResponse** — Odpowiedz mowi, czego rdzen nie rozpoznal — bilans zamiast ciszy

**DesignColorVisionSimulateRequest** — Wynik jest nowym zasobem — zrodlo zostaje nietkniete

**DesignColorVisionSimulateResponse** — Wynik jest nowym zasobem — zrodlo zostaje nietkniete

**DesignColorAccessibilityAuditRequest** — Czynnosc CZYTA i niczego nie zmienia

**DesignColorAccessibilityAuditResponse** — Czynnosc CZYTA i niczego nie zmienia

**DesignStockSearchRequest** — Brak skonfigurowanego dostawcy jest ODPOWIEDZIA nazywajaca brak, nie cisza

**DesignStockSearchResponse** — Brak skonfigurowanego dostawcy jest ODPOWIEDZIA nazywajaca brak, nie cisza

**DesignPrintPreflightRequest** — Czynnosc CZYTA i niczego nie zmienia; zastrzezenie z waga blad zatrzyma druk, wiec Operator ma je zobaczyc TUTAJ, a nie w drukarni

**DesignPrintPreflightResponse** — Czynnosc CZYTA i niczego nie zmienia; zastrzezenie z waga blad zatrzyma druk, wiec Operator ma je zobaczyc TUTAJ, a nie w drukarni

**DesignLargeformatTileRequest** — Odpowiedz niesie uklad kafli, wiec Operator wie, ktory kafel gdzie idzie, zanim cokolwiek wydrukuje

**DesignLargeformatTileResponse** — Odpowiedz niesie uklad kafli, wiec Operator wie, ktory kafel gdzie idzie, zanim cokolwiek wydrukuje

**DesignChartRenderRequest** — Dane przychodza seriami, a nie obrazem — wykres da sie wiec przerysowac po zmianie liczb, zamiast rysowac go od nowa

**DesignChartRenderResponse** — Dane przychodza seriami, a nie obrazem — wykres da sie wiec przerysowac po zmianie liczb, zamiast rysowac go od nowa

**DesignPhotoCropRequest** — Wynik jest wariantem zrodla — oryginal zostaje nietkniety

**DesignPhotoCropResponse** — Wynik jest wariantem zrodla — oryginal zostaje nietkniety

**DesignPhotoUpscaleRequest** — Gdy kanal modelu obrazowego stoi, liczy kanal; gdy nie stoi, liczy rachunek wkompilowany — odpowiedz mowi to polem computedBy

**DesignPhotoUpscaleResponse** — Gdy kanal modelu obrazowego stoi, liczy kanal; gdy nie stoi, liczy rachunek wkompilowany — odpowiedz mowi to polem computedBy

**DesignPhotoInpaintRequest** — Gdy kanal modelu obrazowego stoi, liczy kanal; gdy nie stoi, liczy rachunek wkompilowany — odpowiedz mowi to polem computedBy

**DesignPhotoInpaintResponse** — Gdy kanal modelu obrazowego stoi, liczy kanal; gdy nie stoi, liczy rachunek wkompilowany — odpowiedz mowi to polem computedBy

**DesignPhotoExpandRequest** — Gdy kanal modelu obrazowego stoi, liczy kanal; gdy nie stoi, liczy rachunek wkompilowany — odpowiedz mowi to polem computedBy

**DesignPhotoExpandResponse** — Gdy kanal modelu obrazowego stoi, liczy kanal; gdy nie stoi, liczy rachunek wkompilowany — odpowiedz mowi to polem computedBy

**DesignPhotoBackgroundRemoveRequest** — Gdy kanal modelu obrazowego stoi, liczy kanal; gdy nie stoi, liczy rachunek wkompilowany — odpowiedz mowi to polem computedBy

**DesignPhotoBackgroundRemoveResponse** — Gdy kanal modelu obrazowego stoi, liczy kanal; gdy nie stoi, liczy rachunek wkompilowany — odpowiedz mowi to polem computedBy

**AodMuteGetRequest** — Wyciszenie przeterminowane nie wchodzi do wykazu — Operator nie ma go odklikiwac

**AodMuteGetResponse** — Wyciszenie przeterminowane nie wchodzi do wykazu — Operator nie ma go odklikiwac

**AodMuteSetRequest** — Zniesienie idzie ta sama komenda z polem muted rownym falszowi — odwracalnosc jednym ruchem jest wymogiem, nie wygoda. Odpowiedz oddaje wykaz wyciszen PO zmianie, zeby powloka nie musiala pytac drugi raz. Zmiana rozglasza sie zdarzeniem aod.mute.changed na pozostale powloki Operatora

**AodMuteSetRequest.MuteId** — Podane przy muted rownym falszowi znosi dokladnie to wyciszenie; puste znosi wyciszenie zlozone z rodzaju i zakresu tego zadania

**AodMuteSetRequest.EndsAt** — Wyciszenie czasowe bez niej jest zadaniem niezgodnym z kontraktem — cisza bez konca nie mowi Operatorowi, do kiedy trwa

**AodMuteSetResponse** — Zniesienie idzie ta sama komenda z polem muted rownym falszowi — odwracalnosc jednym ruchem jest wymogiem, nie wygoda. Odpowiedz oddaje wykaz wyciszen PO zmianie, zeby powloka nie musiala pytac drugi raz. Zmiana rozglasza sie zdarzeniem aod.mute.changed na pozostale powloki Operatora

**AodMuteSetResponse.Changed** — Falsz znaczy wyciszenie powtorzone albo zniesienie wyciszenia, ktorego nie bylo

**AodSignalReportRequest** — Nosnik dla trzech klas rozdz. 3.2, ktorych rdzen nie widzi wlasna telemetria: wyniku kontroli jakosci, harmonogramu przebiegow automatyk i powtarzalnosci czynnosci Operatora. Odpowiedz mowi wprost, czy sygnal wpadl w wyciszenie i ktore — sygnal wyciszony odklada sie nadal, bo wyciszenie wstrzymuje UJAWNIENIE, a nie zapis

**AodSignalReportResponse** — Nosnik dla trzech klas rozdz. 3.2, ktorych rdzen nie widzi wlasna telemetria: wyniku kontroli jakosci, harmonogramu przebiegow automatyk i powtarzalnosci czynnosci Operatora. Odpowiedz mowi wprost, czy sygnal wpadl w wyciszenie i ktore — sygnal wyciszony odklada sie nadal, bo wyciszenie wstrzymuje UJAWNIENIE, a nie zapis

**AodSignalListRequest** — Odpowiedz nazywa wyciszenia, ktore sygnaly wstrzymaly, i liczbe wstrzymanych — cisza, po ktorej Operator nie wie, ze cos milczy, jest gorsza od braku wyciszenia

**AodSignalListResponse** — Odpowiedz nazywa wyciszenia, ktore sygnaly wstrzymaly, i liczbe wstrzymanych — cisza, po ktorej Operator nie wie, ze cos milczy, jest gorsza od braku wyciszenia

**MemoryDisableListRequest** — Odczyt do pary z memory.disable.set; okno konfiguracji buduje z niego zakres ustawien „Pamiec"

**MemoryDisableListResponse** — Odczyt do pary z memory.disable.set; okno konfiguracji buduje z niego zakres ustawien „Pamiec"

**MemoryDisableSetRequest** — Zniesienie idzie ta sama komenda z polem disabled rownym falszowi: odwracalnosc jednym ruchem. WYLACZENIE NIE KASUJE TRESCI — tym rozni sie od memory.delete; wpis wylaczony zostaje na miejscu i wraca w calosci po zniesieniu. Odpowiedz oddaje wykaz wylaczen PO zmianie

**MemoryDisableSetResponse** — Zniesienie idzie ta sama komenda z polem disabled rownym falszowi: odwracalnosc jednym ruchem. WYLACZENIE NIE KASUJE TRESCI — tym rozni sie od memory.delete; wpis wylaczony zostaje na miejscu i wraca w calosci po zniesieniu. Odpowiedz oddaje wykaz wylaczen PO zmianie

**MemoryDisableSetResponse.Changed** — Falsz znaczy wylaczenie powtorzone albo zniesienie wylaczenia, ktorego nie bylo

**NotificationListRequest** — Filtry zawezaja wykaz do klas, wag, stanow albo zrodla; pominiete zwracaja rejestr niezamkniety. Licznik zdarzen nowych idzie zawsze, niezaleznie od filtru — plakietka paska kontekstu liczy caly rejestr, nie widok

**NotificationListResponse** — Filtry zawezaja wykaz do klas, wag, stanow albo zrodla; pominiete zwracaja rejestr niezamkniety. Licznik zdarzen nowych idzie zawsze, niezaleznie od filtru — plakietka paska kontekstu liczy caly rejestr, nie widok

**NotificationAcknowledgeRequest** — Wykaz pusty oznacza wszystkie zdarzenia nowe, czyli dzialanie zbiorcze centrum

**NotificationAcknowledgeResponse** — Wykaz pusty oznacza wszystkie zdarzenia nowe, czyli dzialanie zbiorcze centrum

**NotificationResolveRequest** — Dzialanie, ktore je zamyka, wykonuje rodzina wlasciwa — ta komenda zapisuje wylacznie skutek w rejestrze centrum

**NotificationResolveResponse** — Dzialanie, ktore je zamyka, wykonuje rodzina wlasciwa — ta komenda zapisuje wylacznie skutek w rejestrze centrum

**NotificationSnoozeRequest** — Po niej zdarzenie wraca do stanu nowe i znow liczy sie do plakietki — odlozenie nie jest zamknieciem

**NotificationSnoozeResponse** — Po niej zdarzenie wraca do stanu nowe i znow liczy sie do plakietki — odlozenie nie jest zamknieciem

**StreamChunkEvent.MessageId** — Znaczenie pola jest STALE od pierwszego fragmentu do domykajacego — nadawca nie ma prawa zmienic go w polowie strumienia. W debacie jest to kod WYPOWIEDZI, nie kod uczestnika; mowce klient bierze z RoundtableStatement.participantId

**AuthChangedEvent** — Sekcja Uwierzytelnianie Okna Ustawien odswieza sie bez odpytywania

**DeviceChangedEvent** — Rozglaszane do WSZYSTKICH polaczonych urzadzen, zeby uniewaznienie jednego bylo natychmiast widoczne na pozostalych ekranach — tej samej drogi uzywa odzyskanie konta, ktore uniewaznia tokeny wydane przed zmiana hasla

**AdvisorConsultedEvent** — Zdarzenie istnieje dla JAWNOSCI: rada, ktorej nie widac w strumieniu, jest przeslanka, ktorej w aktach nie ma — a agent na niej dziala. Niesie okno, doradce, podstawe doboru i skrot rady; pelna tresc wraca wynikiem komendy advisor.consult i zapisem dziennika konsultacji

**AppsWorkspaceChangedEvent** — Bez tego zdarzenia apps.workspace.update byl jedyna komenda zapisu modulu, ktora nic nie rozglaszala — drugie okno nie dowiadywalo sie o zmianie

**SessionToolAttachedEvent** — Skoro model dostaje narzedzie, ktorego nie mial, Operator MUSI to zobaczyc — kazde posuniecie widac na ekranie

**SessionToolDetachedEvent** — Widocznosc obowiazuje w OBIE strony: zestaw, ktory urosl na oczach Operatora, nie moze skurczyc sie po cichu — bez tego zdarzenia sasiednie okno pokazywaloby narzedzie, ktorego model juz nie ma

**AppsPreviewChangedEvent** — Bez tego zdarzenia podglad na zywo bylby odswiezaniem na zadanie, czyli nie podgladem na zywo

**AppsDeploymentLogEvent** — Osobne od apps.build.changed, ktore niesie STAN wdrozenia, a nie jego dziennik

**SpeechListenPartialEvent** — Tekst jest nietrwaly i bywa poprawiany kolejnym zdarzeniem — trwaly zapis daje dopiero speech.transcribe

**AutomationExecutionLoggedEvent** — Nazwa rozna od komendy automation.execution.log, bo kontrakt nie dopuszcza tej samej nazwy w obu wykazach

**DesignBoardChangedEvent** — Dzis kompozycja jezdzi w obie strony komendami, ale drugie okno i drugie polaczenie nie dowiaduja sie o zmianie niczym — zdarzenia obszaru jest dokladnie jedno i dotyczy zasobu

**DesignBoardPresenceEvent** — Zdarzenie ULOTNE — nie ma odpowiednika w bazie i nie jest odtwarzane po ponownym polaczeniu

**AlertTriggeredEvent** — Zdarzenie jest droga alertu do Always On Display i do powiadomienia w aplikacji; rejestr odpytywany komenda alert.trigger.list sam by tam nie dotarl, bo alert ma znaczenie wtedy, gdy Operator nie patrzy

**AodMuteChangedEvent** — Rozgloszenie sciga cisza pozostale powloki Operatora: bez niego Operator wyciszalby w jednej, a sugestie wchodzilyby w drugiej

**MemoryDisableChangedEvent** — Rozgloszenie idzie na pozostale powloki Operatora, zeby wykaz pamieci w kazdej z nich pokazywal ten sam stan wlaczenia

**NotificationRaisedEvent** — Rozglaszane do wszystkich polaczen Operatora, zeby plakietka i kolumna centrum nie musialy odpytywac rdzenia w tle; to samo zdarzenie jest podstawa Toastu w chwili wystapienia

**NotificationChangedEvent** — Rozglaszane do wszystkich polaczen, zeby obsluga na jednym urzadzeniu byla natychmiast widoczna na pozostalych

**UnknownCommandPayload** — Fail-open: polaczenie nie jest zrywane, sesja nie jest blokowana, kolejne zadania sa przyjmowane.

**ZdarzenieNieznanej** — Fail-open: polaczenie nie jest zrywane, a obszar spoza kontraktu dostaje zdarzenie zapasowe.

**NarzedziaModelu** — Kanal modelu dostaje sterowanie platforma jako narzedzia, a nie jako osobny parser intencji; kazda deklaracja odpowiada komendzie kontraktu.

**KnownModuleIds** — Lista informacyjna, nie brama.

**KnownChannelKinds** — Lista informacyjna; nowy kanal = nowy wiersz rejestru, nie zmiana kodu.

**KnownSettingKeys** — Lista informacyjna, nie brama — katalog ustawien jest sterowany danymi, wiec nowa pozycja to nowy wiersz, nie zmiana kodu.

**KnownSessionConfigKeys** — Klucz powstaje mechanicznie: "sesja.konfiguracja." + wartosc SessionConfigArea. Konfiguracja sesji nie ma wlasnego rezolwera — idzie tym samym rezolwerem osmiu poziomow i trzech osi, ktory obsluguje config.get i config.set. Lista informacyjna, nie brama.

### budowa/server/internal/dane/studio_kontrola_pracy.go

Uzasadnienia i zastrzeżenia projektowe przeniesione z komentarzy warstwy danych blokad, dziennika czynności i kopii zapasowych modułu Studio.

**Plik** — Trzy obszary (blokady fragmentów, dziennik czynności wraz z zależnościami, kopie zapasowe i nastawy pracy) leżą w jednym pliku, bo wiąże je jedno pytanie zadawane w jednym miejscu: czy tę zmianę wolno wnieść, a jeśli tak, to czym ją potem cofnąć. Blokada odpowiada na pierwszą połowę, dziennik na drugą, a kopia zapasowa jest siatką pod obiema — zakłada się ją przed czynnością nieodwracalną i po nieudanym zapisie. Rozdzielenie ich na trzy pliki rozdzieliłoby zapytania, które i tak padają razem. Znakowanie fragmentów, zajęcia wykonawców i spięcia leżą osobno, w warstwie danych znakowania wykonawcy: tamte opisują, co ktoś o dokumencie powiedział i kto nad nim pracuje, a nie czego nie wolno tknąć. Wygasanie zajęcia i wygasanie kopii zapasowej liczy się w SQL (`strftime('now')`), bo obie wielkości muszą być tym samym zegarem co kolumny `utworzono` i `wygasa` w tych wierszach — zegar rdzenia i zegar bazy rozjadą się przy pierwszej różnicy strefy, a wtedy kopia wygasłaby wcześniej albo później, niż mówi nastawa Operatora.

**BlokadaFragmentuStudia** — Pole `Zasieg` rozstrzyga, kogo blokada dotyczy: wartość `model` (postać domyślna) wiąże wyłącznie wykonawców, wartość `everyone` wiąże także Operatora. Operator zmienia fragment zablokowany bez przeszkód, dopóki nie zażąda zasięgu `everyone` jawnie.

**CzynnoscDokumentuStudia** — Pola `StanPrzed` i `StanPo` niosą wycinek objęty czynnością, nie migawkę całego dokumentu — inaczej cofnięcie czynności ze środka dziennika zabrałoby ze sobą wszystko, co po niej weszło, czyli byłoby przywróceniem wersji zamiast cofnięciem jednej zmiany.

**KopiaZapasowaStudia** — Pola `UdaloSie` i `PowodNiepowodzenia` istnieją, bo wskaźnik zapisano pokazany przy zapisie nieudanym jest najgorszym możliwym błędem tego modułu: Operator zamknie okno i straci pracę. Kopia nieudana zostaje wierszem — musi być widoczna i nazwana, a nie zniknąć razem z niepowodzeniem.

**blokadyStudiaPrzesun** — Przesunięcie zakresu po wpisie następuje przed blokadą leżącą za punktem edycji. Bez tego blokada zaczęłaby po pierwszej edycji chronić nie ten fragment, co miała, i to bez śladu. Kolumna `przesuniecia` rośnie razem z zakresem: jest miarą zaufania do zakresu, widoczną dla Operatora w wykazie.

**kopieStudiaPrzemiec** — Zasada wygasania kopii jest jawnym, odwracalnym ustawieniem Operatora, więc granice podaje wołający, a nie stała rdzenia. Kopia nieudana nie wygasa razem z udanymi: jest jedynym śladem, że praca nie doszła na dysk, i ma zostać, dopóki Operator jej nie zobaczy.

**PrzesunBlokadyFragmentow** — Blokada leżąca przed punktem edycji zostaje nietknięta; leżąca za nim jedzie o różnicę długości. Blokada, w której środek trafiła edycja, też zostaje nietknięta — skoro edycja weszła w blokadę, wolno jej było tam wejść (Operator albo zasięg `model` przy czynności Operatora), a wtedy zakres blokady ma zostać taki, jaki Operator ustawił.

**ZapiszCzynnoscDokumentu** — Kolejność nadaje baza, nie wołający: dwie czynności zapisane w tej samej chwili muszą dostać różne numery, a numer nadany przez rdzeń z odczytu stanu bieżącego byłby wyścigiem. Więz UNIQUE(dokument, kolejnosc) zamienia ten wyścig w błąd zapisu, zamiast w dwa wpisy o tym samym miejscu w porządku.

**CzynnosciDokumentu** — Zależności doczytywane są jednym zapytaniem dla całego dokumentu, nie zapytaniem na wpis: dziennik długiego dokumentu ma setki wpisów, a pytanie na każdy wpis zamieniłoby odczyt wykazu w setki zapytań.

**NastawaPracy** — Wiersz nastaw zakłada się przy odczycie, nie przy zapisie, bo odczyt nastaw autozapisu ma oddać nastawy obowiązujące także wtedy, gdy Operator nigdy ich nie ruszał — a wtedy obowiązują wartości domyślne kolumn, i to jest odpowiedź, nie brak danych.

**ZapiszSkutekAutozapisu** — Skutek zapisu odkłada się osobnym poleceniem, bo pisze go inna czynność niż nastawy: nastawy stawia Operator, skutek zapisuje sam mechanizm autozapisu. Jedno wspólne polecenie kazałoby autozapisowi przepisywać nastawy, których nie zmieniał.

### budowa/server/internal/dane/agent_zakres.go

Uzasadnienia i zastrzeżenia projektowe przeniesione z komentarzy zakresu działania eksperta widzianego od strony eksperta.

**Plik** — Repozytorium osobne, a nie kolejne czynności `RepozytoriumAgentow`: tamten kontrakt opisuje bibliotekę ekspertów — założenie, wykaz, zmianę, usunięcie — i pracuje w nim równolegle więcej rąk. Zakres działania jest inną odpowiedzialnością: odpowiada nie na pytanie „jaki jest ten ekspert", lecz „co temu ekspertowi wolno zrobić w systemie". Zapis, którego nikt nie czyta przy wykonaniu, jest gorszy niż jego brak — wiersze zakładane tutaj czyta straż zakresu eksperta (`core/straz_eksperta.go`) i na ich podstawie odmawia: nałożenia eksperta na okno modułu, którego nie ma w jego zakresie, oraz powołania podagentów przez eksperta z wyłączonym Subagent Network. Bez tego Operator widziałby ograniczenie, którego nikt nie egzekwuje. Brak wiersza znaczy stan wyjściowy platformy, czyli pełny dostęp: brak modułów znaczy „wszędzie", brak przełącznika izolacji znaczy „nieodcięty".

**PrzypisaniaEkspertow** — Dwa zapytania, nie jedno z UNION: źródła mają inne kolumny i inne znaczenie celu — projekt ma kod własny, stanowisko obsady ma klucz wiersza — a złożenie ich w jedno zapytanie kazałoby czytelnikowi zgadywać, skąd wziął się wiersz.

### budowa/server/internal/dane/studio_znakowanie_wykonawcy.go

Uzasadnienia i zastrzeżenia projektowe przeniesione z komentarzy części drugiej warstwy danych kontroli pracy modułu Studio.

**Plik, wersje w szeregach** — `studio_wersje.go` zna wersję sprzed dobudowy: treść, etykietę, kamień milowy. Kolumny `szereg`, `postac_json` i tożsamości wykonawcy dołożyła migracja 367 i pyta o nie wyłącznie ten odcinek — wykaz historii z przełącznikiem „pokaż także zapisy samoczynne" oraz porównanie postaci dwóch wersji. Dopisanie ich do tamtego pliku byłoby wejściem w plik cudzego odcinka; osobny odczyt tych samych wierszy nie zakłada drugiego pojęcia wersji, bo tabela jest jedna i przywraca się ją tą samą drogą.

**Plik, zajęcie a blokada** — Blokada Operatora jest trwała i skierowana przeciw wykonawcom: zdejmuje ją wyłącznie Operator. Zajęcie fragmentu jest chwilowe, wygasa samo i chroni przed drugim wykonawcą. Dwa różne byty, dwie tabele — zlanie ich dałoby blokadę, która wygasa (czyli żadną), albo zajęcie, którego nikt nie zdejmie po agencie ubitym w pół pracy.

**ZajecieFragmentuStudia.Wygasle** — Liczone zapytaniem SQL, nie w rdzeniu: zegar rdzenia i zegar bazy rozjadą się przy pierwszej różnicy strefy, a wtedy zajęcie trzymałoby fragment dłużej albo krócej, niż mówi jego własna kolumna `wygasa`.

**wersjaZalozycielskaStudia** — Szereg autozapisu odpada z rachunku z zamysłem: gdyby pierwszy zapis samoczynny wypadł przed pierwszym zapisem Operatora, powrót do stanu pierwotnego wracałby do stanu przypadkowego, a nie do tego, co Operator założył.

**ZapiszZajecieFragmentu** — Odstęp zerowy albo ujemny znaczy zajęcie bez wygasania — i takiego zajęcia rdzeń nie zakłada, bo agent ubity w pół pracy trzymałby fragment na zawsze; granicę podaje wołający.

**ZapiszWersjeSzeregu** — Wskaźnika wersji bieżącej nie przestawia, tak samo jak `ZapiszWersje` z `studio_wersje.go`; decyzja, kiedy nowa wersja staje się bieżącą, należy do wołającego.

### budowa/server/internal/dane/badania_dobudowa.go

Uzasadnienia i zastrzeżenia projektowe przeniesione z komentarzy dobudowy obszaru Research po stronie danych.

**Plik** — Ustalenia i ich otoczenie leżą w `badania_dobudowa_ustalenia.go`, odkrywanie i przestrzeń w `badania_dobudowa_odkrycia.go`, raport i eksport w `badania_dobudowa_raport.go`. Wszystkie cztery pliki niosą metody tego samego `*repozytoriumBadan`, tak jak trzy pliki zastane. Kontrakt dobudowy jest jednym interfejsem wpiętym do `RepozytoriumBadan` przez osadzenie: jedna deklaracja, jedno miejsce do przeczytania, a plik `badania.go` nie rośnie o siedemdziesiąt sygnatur, których nie implementuje.

## budowa/server/internal/dane/urzadzenia_powiadomien.go

Rozbicie tabel `urzadzenie_powiadomien`, `powiadomienie` i `powiadomienie_dostarczenie`
na trzy repozytoria zmusiłoby silnik do składania transakcji z kawałków trzech
właścicieli, a najważniejsza droga tego pliku (Takt) jest transakcją obejmującą
wszystkie trzy tabele naraz.

Konstruktor `NowePowiadomienia` jest eksportowany i bierze uchwyt puli, żeby
pakiet silnika powiadomień mógł złożyć repozytorium z uchwytu, który już dostał
w kompozycji, bez dodawania pola w `dane/zestaw.go`. Pula połączeń jest jedna na
proces; to repozytorium jej nie otwiera i nie zamyka.

`KluczKanalu` w `RejestracjaPowiadomien` dla kanału `polaczenie` niesie ten sam
napis, który transport przekazuje jako `Tozsamosc.IdKlienta` — adres gniazda,
nie odrębną tożsamość urządzenia.

`Wyslij` w `PlanTaktu` jest domknięciem, a nie zwracaną wartością, żeby sprawdzenie
stanu wiersza i sama wysyłka zmieściły się w jednej transakcji zapisu: sprawdzenie
przy pobraniu należnych nie wystarcza, bo między pobraniem a wysłaniem decyzja
może zapaść i urządzenie dostałoby powiadomienie po fakcie. Pusta lista przyjętych
odbiorców znaczy „nie było komu" i jest powodem ponowienia, nie błędem.

Kolejność kroków w `Takt` nie może zostać przestawiona: `zajmijPowiadomienie`
bierze zamek zapisu i w tym samym poleceniu sprawdza, że wiersz nadal jest
`oczekuje` — zero zmienionych wierszy znaczy, że decyzja już zapadła, termin minął
albo doręczenie już było. Dopiero pod tym zamkiem czytany jest wiersz i sprawdzany
termin ważności. Wysyłka idzie wewnątrz transakcji, więc koperta wychodzi przy
trzymanym zamku zapisu — cena możliwa do przyjęcia, bo doręczenie kanałem
`polaczenie` jest zapisem do gniazda w pamięci procesu. Odwrotna kolejność —
zwolnić zamek, potem wysłać — przywraca szczelinę, w której powiadomienie
o zapadłej już decyzji zdąży wyjść. Wołający, który nie odda `Wyslij` albo
`NastepnaProba`, dostaje błąd zamiast cichego pominięcia.

Zapytanie `Nalezne` porządkuje wynik tak, że bez tego kolejka pilna stałaby za
zwykłą, która akurat weszła pierwsza, i priorytet nie znaczyłby nic.
