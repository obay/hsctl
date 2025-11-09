package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/obay/hsctl/internal/hubspot"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var contactsCmd = &cobra.Command{
	Use:   "contacts",
	Short: "Manage HubSpot contacts",
	Long:  `Manage HubSpot contacts with CRUD operations.`,
}

var listContactsCmd = &cobra.Command{
	Use:   "list",
	Short: "List all contacts",
	Long:  `List all contacts in HubSpot with their properties and email addresses.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		apiKey := getAPIKey()
		if apiKey == "" {
			return fmt.Errorf("API key is required. Set HUBSPOT_API_KEY env var or use --api-key flag")
		}

		client := hubspot.NewClient(apiKey)
		limit, _ := cmd.Flags().GetInt("limit")
		if limit == 0 {
			limit = 100
		}

		format, _ := cmd.Flags().GetString("format")
		showAll, _ := cmd.Flags().GetBool("all")

		var allContacts []hubspot.Contact
		after := ""

		for {
			contactResp, err := client.ListContacts(limit, after)
			if err != nil {
				return fmt.Errorf("failed to list contacts: %w", err)
			}

			allContacts = append(allContacts, contactResp.Results...)

			if !showAll || contactResp.Paging == nil || contactResp.Paging.Next == nil {
				break
			}
			after = contactResp.Paging.Next.After
		}

		return printContacts(allContacts, format)
	},
}

var listPropertiesCmd = &cobra.Command{
	Use:   "properties",
	Short: "List all contact properties",
	Long:  `List all available properties for contacts in HubSpot.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		apiKey := getAPIKey()
		if apiKey == "" {
			return fmt.Errorf("API key is required. Set HUBSPOT_API_KEY env var or use --api-key flag")
		}

		client := hubspot.NewClient(apiKey)
		properties, err := client.ListProperties()
		if err != nil {
			return fmt.Errorf("failed to list properties: %w", err)
		}

		format, _ := cmd.Flags().GetString("format")
		return printProperties(properties, format)
	},
}

var createContactCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new contact",
	Long:  `Create a new contact in HubSpot.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		apiKey := getAPIKey()
		if apiKey == "" {
			return fmt.Errorf("API key is required. Set HUBSPOT_API_KEY env var or use --api-key flag")
		}

		client := hubspot.NewClient(apiKey)

		properties := make(map[string]interface{})
		email, _ := cmd.Flags().GetString("email")
		firstName, _ := cmd.Flags().GetString("firstname")
		lastName, _ := cmd.Flags().GetString("lastname")
		lifecycleStage, _ := cmd.Flags().GetString("lifecycle-stage")
		propertiesStr, _ := cmd.Flags().GetString("properties")

		if email != "" {
			properties["email"] = email
		}
		if firstName != "" {
			properties["firstname"] = firstName
		}
		if lastName != "" {
			properties["lastname"] = lastName
		}
		if lifecycleStage != "" {
			properties["lifecyclestage"] = lifecycleStage
		}

		// Parse additional properties from string (format: "key1=value1,key2=value2")
		if propertiesStr != "" {
			pairs := strings.Split(propertiesStr, ",")
			for _, pair := range pairs {
				parts := strings.SplitN(pair, "=", 2)
				if len(parts) == 2 {
					key := strings.TrimSpace(parts[0])
					value := strings.TrimSpace(parts[1])
					properties[key] = value
				}
			}
		}

		if len(properties) == 0 {
			return fmt.Errorf("at least one property is required to create a contact")
		}

		contact, err := client.CreateContact(properties)
		if err != nil {
			return fmt.Errorf("failed to create contact: %w", err)
		}

		fmt.Println("Contact created successfully:")
		return printContacts([]hubspot.Contact{*contact}, "table")
	},
}

var updateContactCmd = &cobra.Command{
	Use:   "update [contact-id]",
	Short: "Update a contact",
	Long:  `Update a contact's properties or lifecycle stage.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		apiKey := getAPIKey()
		if apiKey == "" {
			return fmt.Errorf("API key is required. Set HUBSPOT_API_KEY env var or use --api-key flag")
		}

		client := hubspot.NewClient(apiKey)
		contactID := args[0]

		properties := make(map[string]interface{})
		email, _ := cmd.Flags().GetString("email")
		firstName, _ := cmd.Flags().GetString("firstname")
		lastName, _ := cmd.Flags().GetString("lastname")
		lifecycleStage, _ := cmd.Flags().GetString("lifecycle-stage")
		propertiesStr, _ := cmd.Flags().GetString("properties")

		if email != "" {
			properties["email"] = email
		}
		if firstName != "" {
			properties["firstname"] = firstName
		}
		if lastName != "" {
			properties["lastname"] = lastName
		}
		if lifecycleStage != "" {
			properties["lifecyclestage"] = lifecycleStage
		}

		// Parse additional properties from string
		if propertiesStr != "" {
			pairs := strings.Split(propertiesStr, ",")
			for _, pair := range pairs {
				parts := strings.SplitN(pair, "=", 2)
				if len(parts) == 2 {
					key := strings.TrimSpace(parts[0])
					value := strings.TrimSpace(parts[1])
					properties[key] = value
				}
			}
		}

		if len(properties) == 0 {
			return fmt.Errorf("at least one property is required to update a contact")
		}

		contact, err := client.UpdateContact(contactID, properties)
		if err != nil {
			return fmt.Errorf("failed to update contact: %w", err)
		}

		fmt.Println("Contact updated successfully:")
		return printContacts([]hubspot.Contact{*contact}, "table")
	},
}

var deleteContactCmd = &cobra.Command{
	Use:   "delete [contact-id]",
	Short: "Delete a contact",
	Long:  `Delete a contact from HubSpot by ID.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		apiKey := getAPIKey()
		if apiKey == "" {
			return fmt.Errorf("API key is required. Set HUBSPOT_API_KEY env var or use --api-key flag")
		}

		client := hubspot.NewClient(apiKey)
		contactID := args[0]

		force, _ := cmd.Flags().GetBool("force")
		if !force {
			contact, err := client.GetContact(contactID)
			if err != nil {
				return fmt.Errorf("failed to get contact: %w", err)
			}

			email := "N/A"
			if e, ok := contact.Properties["email"].(string); ok {
				email = e
			}

			fmt.Printf("Are you sure you want to delete contact %s (email: %s)? [y/N]: ", contactID, email)
			var response string
			fmt.Scanln(&response)
			if strings.ToLower(response) != "y" && strings.ToLower(response) != "yes" {
				fmt.Println("Deletion cancelled.")
				return nil
			}
		}

		err := client.DeleteContact(contactID)
		if err != nil {
			return fmt.Errorf("failed to delete contact: %w", err)
		}

		fmt.Printf("Contact %s deleted successfully.\n", contactID)
		return nil
	},
}

var queryContactsCmd = &cobra.Command{
	Use:   "query [search-query]",
	Short: "Search for contacts",
	Long:  `Search for contacts using HubSpot's search API.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		apiKey := getAPIKey()
		if apiKey == "" {
			return fmt.Errorf("API key is required. Set HUBSPOT_API_KEY env var or use --api-key flag")
		}

		client := hubspot.NewClient(apiKey)
		query := args[0]
		limit, _ := cmd.Flags().GetInt("limit")
		if limit == 0 {
			limit = 100
		}

		format, _ := cmd.Flags().GetString("format")

		contactResp, err := client.SearchContacts(query, limit)
		if err != nil {
			return fmt.Errorf("failed to search contacts: %w", err)
		}

		return printContacts(contactResp.Results, format)
	},
}

func init() {
	rootCmd.AddCommand(contactsCmd)

	// List contacts command
	contactsCmd.AddCommand(listContactsCmd)
	listContactsCmd.Flags().IntP("limit", "l", 100, "Maximum number of contacts to retrieve")
	listContactsCmd.Flags().BoolP("all", "a", false, "Retrieve all contacts (paginate through all pages)")
	listContactsCmd.Flags().StringP("format", "f", "table", "Output format (table, json)")

	// List properties command
	contactsCmd.AddCommand(listPropertiesCmd)
	listPropertiesCmd.Flags().StringP("format", "f", "table", "Output format (table, json)")

	// Create contact command
	contactsCmd.AddCommand(createContactCmd)
	createContactCmd.Flags().StringP("email", "e", "", "Email address")
	createContactCmd.Flags().StringP("firstname", "f", "", "First name")
	createContactCmd.Flags().StringP("lastname", "l", "", "Last name")
	createContactCmd.Flags().String("lifecycle-stage", "", "Lifecycle stage (e.g., lead, customer)")
	createContactCmd.Flags().StringP("properties", "p", "", "Additional properties (format: key1=value1,key2=value2)")

	// Update contact command
	contactsCmd.AddCommand(updateContactCmd)
	updateContactCmd.Flags().StringP("email", "e", "", "Email address")
	updateContactCmd.Flags().StringP("firstname", "f", "", "First name")
	updateContactCmd.Flags().StringP("lastname", "l", "", "Last name")
	updateContactCmd.Flags().String("lifecycle-stage", "", "Lifecycle stage (e.g., lead, customer)")
	updateContactCmd.Flags().StringP("properties", "p", "", "Additional properties (format: key1=value1,key2=value2)")

	// Delete contact command
	contactsCmd.AddCommand(deleteContactCmd)
	deleteContactCmd.Flags().Bool("force", false, "Skip confirmation prompt")

	// Query contacts command
	contactsCmd.AddCommand(queryContactsCmd)
	queryContactsCmd.Flags().IntP("limit", "l", 100, "Maximum number of results")
	queryContactsCmd.Flags().StringP("format", "f", "table", "Output format (table, json)")
}

func getAPIKey() string {
	// Check flag first
	if apiKey := viper.GetString("api-key"); apiKey != "" {
		return apiKey
	}
	// Check environment variable
	if apiKey := os.Getenv("HUBSPOT_API_KEY"); apiKey != "" {
		return apiKey
	}
	return ""
}

func printContacts(contacts []hubspot.Contact, format string) error {
	if format == "json" {
		jsonData, err := json.MarshalIndent(contacts, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(jsonData))
		return nil
	}

	// Table format
	fmt.Printf("%-20s %-40s %-20s %-20s %-20s\n", "ID", "Email", "First Name", "Last Name", "Lifecycle Stage")
	fmt.Println(strings.Repeat("-", 120))

	for _, contact := range contacts {
		email := getStringValue(contact.Properties["email"])
		firstName := getStringValue(contact.Properties["firstname"])
		lastName := getStringValue(contact.Properties["lastname"])
		lifecycleStage := getStringValue(contact.Properties["lifecyclestage"])

		fmt.Printf("%-20s %-40s %-20s %-20s %-20s\n",
			contact.ID, email, firstName, lastName, lifecycleStage)
	}

	fmt.Printf("\nTotal: %d contact(s)\n", len(contacts))
	return nil
}

func printProperties(properties []hubspot.Property, format string) error {
	if format == "json" {
		jsonData, err := json.MarshalIndent(properties, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(jsonData))
		return nil
	}

	// Table format
	fmt.Printf("%-30s %-30s %-20s %-15s\n", "Name", "Label", "Type", "Field Type")
	fmt.Println(strings.Repeat("-", 95))

	for _, prop := range properties {
		fmt.Printf("%-30s %-30s %-20s %-15s\n",
			prop.Name, prop.Label, prop.Type, prop.FieldType)
	}

	fmt.Printf("\nTotal: %d property(ies)\n", len(properties))
	return nil
}

func getStringValue(v interface{}) string {
	if v == nil {
		return ""
	}
	switch val := v.(type) {
	case string:
		return val
	case float64:
		return strconv.FormatFloat(val, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(val)
	default:
		return fmt.Sprintf("%v", val)
	}
}
