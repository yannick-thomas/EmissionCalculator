import tkinter as tk
from tkinter import ttk
import customtkinter



class CalcApp(customtkinter.CTk):

    WIDTH = 980
    HEIGHT = 720

    CURRENT_PAGE = ''

    def __init__(self):
        super().__init__()

        self.delivery_quantity_briquettes = tk.StringVar(self)
        self.delivery_quantity_oil = tk.StringVar(self)
        self.result_briquettes = 0
        self.is_valid_calc = False

        # Window configuration
        self.title("Emissionsrechner")
        self.geometry(f"{CalcApp.WIDTH}x{CalcApp.HEIGHT}")
        self.minsize(self.WIDTH, self.HEIGHT)
        self.maxsize(self.WIDTH, self.HEIGHT)

        self.grid_columnconfigure(1, weight=1)
        self.grid_columnconfigure((2, 3), weight=0)
        self.grid_rowconfigure((0, 1, 2), weight=1)

        ## Sidebar
        self.sidebar = customtkinter.CTkFrame(self, width=160, corner_radius=0)
        self.sidebar.grid(row=0, column=0, rowspan=5, sticky="nsew")
        self.sidebar.grid_rowconfigure(5, weight=1)
        self.logo_label = customtkinter.CTkLabel(self.sidebar, text="Emissionsrechner", font=customtkinter.CTkFont(size=25, weight="bold"))
        self.logo_label.grid(row=0, column=0, padx=20, pady=(20, 10))
        self.sidebar_briquettes = customtkinter.CTkButton(self.sidebar, command=self.render_briquettes, text="Briketts")
        self.sidebar_briquettes.grid(row=1, column=0, padx=20, pady=10)
        self.sidebar_heating_oil = customtkinter.CTkButton(self.sidebar, command=self.render_heating_oil, text="Heizöl")
        self.sidebar_heating_oil.grid(row=2, column= 0, padx= 20, pady=10)

        ## Content
        #self.main_content = customtkinter.CTkFrame(self, corner_radius=0, fg_color="#232323")
        #self.main_content.grid(row=0, column=1, rowspan=5, sticky="nsew")
        #self.main_content.grid_rowconfigure(5, weight=1)
        #self.main_content_heading = customtkinter.CTkLabel(self.main_content, text="", font=customtkinter.CTkFont(size=20, weight="bold"))
        #self.main_content_heading.grid(row=1, column=0, padx=(20, 0), pady=(20, 0), columnspan=5 , sticky="nw")



    ## functions
    def render_briquettes(self):
        self.is_valid_calc = False
        if (self.CURRENT_PAGE == "heating_oil"):
            self.main_content.destroy()
        self.main_content = customtkinter.CTkFrame(self, corner_radius=0, fg_color="#232323")
        self.main_content.grid(row=0, column=1, rowspan=5, sticky="nsew")
        self.main_content.grid_rowconfigure(5, weight=1)
        self.main_content_heading = customtkinter.CTkLabel(self.main_content, text="Emissionsberechnung für Briketts", font=customtkinter.CTkFont(size=20, weight="bold"))
        self.main_content_heading.grid(row=1, column=0, padx=(20, 0), pady=(20, 0), columnspan=5 , sticky="nw")
        self.fuel_emission_calc_formula = customtkinter.CTkLabel(self.main_content, text_color="white", font=customtkinter.CTkFont(size=14, weight="bold"))
        self.fuel_emission_calc_equals_result = customtkinter.CTkLabel(self.main_content, font=customtkinter.CTkFont(size=14, weight="bold"))
        self.price_comp_co2_costs_formula = customtkinter.CTkLabel(self.main_content, font=customtkinter.CTkFont(size=14, weight="bold"))
        self.price_comp_co2_costs_label = customtkinter.CTkLabel(self.main_content, font=customtkinter.CTkFont(size=16), text_color="grey")
        self.energy_content_formula = customtkinter.CTkLabel(self.main_content, font=customtkinter.CTkFont(size=14, weight="bold"))
        self.result_energy_content_label = customtkinter.CTkLabel(
            self.main_content,
            font=customtkinter.CTkFont(size=14, weight="bold")
        )
        self.result_co2_costs_label = customtkinter.CTkLabel(
            self.main_content,
            font=customtkinter.CTkFont(size=14, weight="bold")
        )
        self.label_quantity_briquettes= customtkinter.CTkLabel(self.main_content, text="Liefermenge:", font=customtkinter.CTkFont(size=16))
        self.label_quantity_briquettes.grid(row=2, column=0,columnspan=1, padx=(20, 0), pady=(20, 0), sticky="nw")
        self.delivery_quanitity_briquettes_entry = customtkinter.CTkEntry(self.main_content, width=50, font=customtkinter.CTkFont(size=12, weight="bold"), textvariable=self.delivery_quantity_briquettes)
        self.delivery_quanitity_briquettes_entry.grid(row=2, column=1, padx=(155, 0), pady=(20, 0), sticky="nw")
        self.delivery_quanitity_briquettes_measure = customtkinter.CTkLabel(self.main_content, text="t", font=customtkinter.CTkFont(size=16, weight="bold"))
        self.delivery_quanitity_briquettes_measure.grid(row=2, column=2, padx=(0, 0), pady=(20, 0), sticky="nw")
        self.calc_button_briquettes = customtkinter.CTkButton(self.main_content, text="Berechnen", width=120, command=self.calc_briquettes)
        self.calc_button_briquettes.grid(row=2, column=5, padx=(20, 0), pady=(20, 0))
        self.label_calc_failure = customtkinter.CTkLabel(self.main_content, font=customtkinter.CTkFont(size=16, weight="bold"))
        self.bind('<Return>', lambda event: self.calc_briquettes())
        self.CURRENT_PAGE = "briquettes"

    def render_heating_oil(self):
        self.is_valid_calc = False
        if (self.CURRENT_PAGE == "briquettes"):
            self.main_content.destroy()
        self.main_content = customtkinter.CTkFrame(self, corner_radius=0, fg_color="#232323")
        self.main_content.grid(row=0, column=1, rowspan=5, sticky="nsew")
        self.main_content.grid_rowconfigure(5, weight=1)
        self.main_content_heading = customtkinter.CTkLabel(self.main_content, text="Emissionsberechnung für Heizöl", font=customtkinter.CTkFont(size=20, weight="bold"))
        self.main_content_heading.grid(row=1, column=0, padx=(20, 0), pady=(20, 0), columnspan=4, sticky="nw")
        self.label_quantity_oil= customtkinter.CTkLabel(self.main_content, text="Liefermenge:", font=customtkinter.CTkFont(size=16))
        self.label_quantity_oil.grid(row=2, column=0,columnspan=1, padx=(20, 0), pady=(20, 0), sticky="nw")
        self.entry_quantity_oil = customtkinter.CTkEntry(self.main_content, width=50, font=customtkinter.CTkFont(size=12, weight="bold"), textvariable=self.delivery_quantity_oil)
        self.entry_quantity_oil.grid(row=2, column=1, padx=(155, 0), pady=(20, 0), sticky="nw")
        self.label_quantity_measure = customtkinter.CTkLabel(self.main_content, text="l", font=customtkinter.CTkFont(size=16, weight="bold"))
        self.label_quantity_measure.grid(row=2, column=2, padx=(2, 0), pady=(20, 0), sticky="nw")
        self.calc_button_oil = customtkinter.CTkButton(self.main_content, text="Berechnen", width=120, command=self.calc_heating_oil)
        self.calc_button_oil.grid(row=2, column=5, padx=(20, 0), pady=(20, 0))
        self.emissions_heading_label = customtkinter.CTkLabel(
            self.main_content,
            font=customtkinter.CTkFont(size=16),
            text_color="grey"
        )
        self.emissions_formula = customtkinter.CTkLabel(
            self.main_content,
            font=customtkinter.CTkFont(size=14, weight="bold")
        )
        self.label_calc_failure = customtkinter.CTkLabel(self.main_content, font=customtkinter.CTkFont(size=16, weight="bold"))
        self.bind('<Return>', lambda event: self.calc_heating_oil())
        self.CURRENT_PAGE = "heating_oil"
        

    def calc_briquettes(self):

        if (self.is_valid_calc):
            self.fuel_emission_calc_label.grid_forget()
            self.fuel_emission_calc_formula.grid_forget()
            self.fuel_emission_calc_equals.grid_forget()
            self.fuel_emission_calc_equals_result.grid_forget()
            self.price_comp_co2_costs_formula.grid_forget()
            self.price_comp_co2_costs_equals.grid_forget()
            self.price_comp_co2_costs_label.grid_forget()
            self.result_co2_costs_label.grid_forget()
            self.energy_content_formula.grid_forget()
            self.energy_content_equals.grid_forget()
            self.energy_content_label.grid_forget()
            self.result_energy_content_label.grid_forget()

        if (not self.is_valid_calc):
            self.label_calc_failure.grid_forget()

        try:
            self.result_briquettes_emissions = 19 * 0.0992 * 1 * float(self.delivery_quantity_briquettes.get().replace(",", "."))
            self.is_valid_calc = True
        except:
            self.label_calc_failure.configure(text="Fehler")
            self.label_calc_failure.grid(row=3, column=1, padx=(100, 0), pady=(30, 0), sticky="nw")
            self.is_valid_calc = False
            return
        
        self.fuel_emission_calc_label = customtkinter.CTkLabel(self.main_content, text="Brennstoffemissionen der aktuellen Lieferung:", font=customtkinter.CTkFont(size=16), text_color="grey")
        self.fuel_emission_calc_label.grid(row=3,columnspan=3, column=0 ,padx=(20, 0), pady=(30, 0), sticky="nw")
        self.fuel_emission_calc_formula = customtkinter.CTkLabel(self.main_content, text="19,0 GJ/t x 0,0992t CO2 / GJ x 1 x " + self.delivery_quantity_briquettes.get() + "t", text_color="white", font=customtkinter.CTkFont(size=14, weight="bold"))
        self.fuel_emission_calc_formula.grid(row=4,columnspan=3, column=0, padx=(20, 0), pady=(0, 0), sticky="nw")
        self.fuel_emission_calc_equals = customtkinter.CTkLabel(self.main_content, text="≈", font=customtkinter.CTkFont(size=18))
        self.fuel_emission_calc_equals.grid(row=4,columnspan=3, column=0, padx=(330, 0), pady=(0, 0), sticky="nw")
        self.fuel_emission_calc_equals_result = customtkinter.CTkLabel(
            self.main_content,
            text=str(round(self.result_briquettes_emissions, 2)).replace(".", ",") + "t CO2",
            font=customtkinter.CTkFont(size=14, weight="bold")
        )
        self.fuel_emission_calc_equals_result.grid(row=4,columnspan=3, column=3, padx=(20, 0), pady=(0, 0), sticky="nw")
        self.price_comp_co2_costs_label = customtkinter.CTkLabel(self.main_content, text="Preisbestandteil CO2 Kosten:", font=customtkinter.CTkFont(size=16), text_color="grey")
        self.price_comp_co2_costs_label.grid(row=5, columnspan=3 ,column=0,padx=(20, 0), pady=(10, 0), sticky="nw")
        self.price_comp_co2_costs_formula = customtkinter.CTkLabel(self.main_content, text=str(self.result_briquettes_emissions).replace(".", ",") + "Kg CO2 x 30€/ CO2 t x 1,19", font=customtkinter.CTkFont(size=14, weight="bold"))
        self.price_comp_co2_costs_formula.grid(row=6, columnspan=3 ,column=0,padx=(20, 0), pady=(0, 470), sticky="nw")
        self.price_comp_co2_costs_equals = customtkinter.CTkLabel(self.main_content, text="=", font=customtkinter.CTkFont(size=18))
        self.price_comp_co2_costs_equals.grid(row=6,columnspan=3, column=0, padx=(330, 0), pady=(0, 0), sticky="nw")
        self.result_co2_costs = self.result_briquettes_emissions * 30 * 1.19
        self.result_co2_costs_label = customtkinter.CTkLabel(
            self.main_content,
            text=str(round(self.result_co2_costs, 2)).replace(".", ",") + " €",
            font=customtkinter.CTkFont(size=14, weight="bold")
        )
        self.result_co2_costs_label.grid(row=6,columnspan=3, column=3, padx=(20, 0), pady=(0, 0), sticky="nw")
        self.energy_content_label = customtkinter.CTkLabel(self.main_content, text="Energiegehalt der Lieferung:", font=customtkinter.CTkFont(size=16), text_color="grey")
        self.energy_content_label.grid(row=6,columnspan=3, column=0, padx=(20, 0), pady=(40, 0), sticky="nw")
        self.energy_content_formula = customtkinter.CTkLabel(self.main_content, text="19,0 GJ x 1 x " + self.delivery_quantity_briquettes.get(), font=customtkinter.CTkFont(size=14, weight="bold"))
        self.energy_content_formula.grid(row=6,columnspan=3, column=0, padx=(20, 0), pady=(65, 0), sticky="nw")
        self.energy_content_equals = customtkinter.CTkLabel(self.main_content, text="=", font=customtkinter.CTkFont(size=18))
        self.energy_content_equals.grid(row=6,columnspan=3, column=0, padx=(330, 0), pady=(65, 0), sticky="nw")
        self.result_energy_content = 19 * 1 * float(self.delivery_quantity_briquettes.get().replace(",", "."))
        self.result_energy_content_label = customtkinter.CTkLabel(
            self.main_content,
            text=str(self.result_energy_content).replace(".", ",") + " GJ",
            font=customtkinter.CTkFont(size=14, weight="bold")
        )
        self.result_energy_content_label.grid(row=6,columnspan=3, column=3, padx=(20, 0), pady=(65, 0), sticky="nw")
        
    def calc_heating_oil(self):
        if (self.is_valid_calc):
            self.emissions_heading_label.grid_forget()
            self.emissions_formula.grid_forget()
            self.emissions_equals.grid_forget()
            self.emissions_result.grid_forget()
            self.emission_component_heading.grid_forget()
            self.emission_component_formula.grid_forget()
            self.emission_component_equals.grid_forget()
            self.emission_component_result_label.grid_forget()

        if (not self.is_valid_calc):
            self.label_calc_failure.grid_forget()

        try:
            self.emissions = 42.8 * 0.074 * 0.845 * float(self.delivery_quantity_oil.get().replace(",", "."))
            self.is_valid_calc = True
        except:
            self.label_calc_failure.configure(text="Fehler")
            self.label_calc_failure.grid(row=3, column=1, padx=(100, 0), pady=(30, 0), sticky="nw")
            self.is_valid_calc = False
            return


        self.emissions_heading_label = customtkinter.CTkLabel(
            self.main_content,
            text="Brennstroffemissionen der Heizöllieferung:",
            font=customtkinter.CTkFont(size=16),
            text_color="grey"
        )
        self.emissions_heading_label.grid(row=3,columnspan=3, column=0 ,padx=(20, 0), pady=(30, 0), sticky="nw")
        self.emissions_formula = customtkinter.CTkLabel(
            self.main_content,
            text="42,8 GJ/t x 0,074 t CO2/GJ x 0,845t/1000l x " + self.delivery_quantity_oil.get()+ "l", 
            font=customtkinter.CTkFont(size=14, weight="bold")
        )
        self.emissions_formula.grid(row=4 ,columnspan=6,column=0, padx=(20, 0), pady=(0, 0), sticky="nw")
        self.emissions_equals = customtkinter.CTkLabel(self.main_content, text="=")
        self.emissions_equals.grid(row=4,columnspan=6, column=2, padx=(50, 0), pady=(0, 0), sticky="nw")
        self.emissions_result = customtkinter.CTkLabel(
            self.main_content,
            text=str(round(self.emissions, 2)).replace(".", ",") + "t CO2",
            font=customtkinter.CTkFont(size=14, weight="bold")
        )
        self.emissions_result.grid(row=4,columnspan=6, column=3, padx=(10, 0), pady=(0, 0), sticky="ne")
        self.emission_component_heading = customtkinter.CTkLabel(
            self.main_content, 
            text="Preisbestandteil CO2 Kosten (inkl. Umsatzsteuer):",
            font=customtkinter.CTkFont(size=16),
            text_color="grey"
        )
        self.emission_component_heading.grid(row=5, columnspan= 6, column=0, padx=(20, 0), pady=(20, 0), sticky="nw")
        self.emission_component_formula = customtkinter.CTkLabel(
            self.main_content,
            text= str(round(self.emissions, 2)).replace(".", ",") + "t CO2 x 30€/t CO2 x 1,12 MwSt.",
            font=customtkinter.CTkFont(size=14, weight="bold")
        )
        self.emission_component_formula.grid(row=6, columnspan=6, column=0, padx=(20, 0), pady=(0, 460), sticky="nw")
        self.emission_component_equals = customtkinter.CTkLabel(
            self.main_content,
            text="="
        )
        self.emission_component_equals.grid(row=6,columnspan=6, column=2, padx=(50, 0), pady=(0, 0), sticky="nw")
        self.emissions = round(self.emissions, 2)
        self.emission_component_result = self.emissions * 30 * 1.19
        self.emission_component_result = round(self.emission_component_result, 2)
        self.emission_component_result_label = customtkinter.CTkLabel(
            self.main_content,
            text=str(self.emission_component_result) + "€",
            font=customtkinter.CTkFont(size=14, weight="bold")
        )
        self.emission_component_result_label.grid(row=6, columnspan=6, column=3, padx=(0, 20), pady=(0, 0), sticky="ne")

if __name__ == "__main__":
    app = CalcApp()
    app.mainloop()