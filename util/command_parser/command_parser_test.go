package command_parser_test

import (
	"io/ioutil"

	"code.cloudfoundry.org/cli/command/commandfakes"
	"code.cloudfoundry.org/cli/command/common"
	"code.cloudfoundry.org/cli/util/command_parser"
	"code.cloudfoundry.org/cli/util/configv3"
	"code.cloudfoundry.org/cli/util/ui"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Command 'Parser'", func() {
	var (
		pluginUI *ui.UI
	)
	BeforeEach(func() {
		var err error
		fakeConfig := new(commandfakes.FakeConfig)
		pluginUI, err = ui.NewPluginUI(fakeConfig, ioutil.Discard, ioutil.Discard)
		Expect(err).ToNot(HaveOccurred())
	})

	It("returns an unknown command error", func() {
		parser, err := command_parser.NewCommandParser()
		Expect(err).ToNot(HaveOccurred())
		status := parser.ParseCommandFromArgs(pluginUI, []string{"howdy"})
		Expect(status).To(Equal(-666))
	})

	Describe("the verbose flag", func() {
		var parser command_parser.CommandParser

		BeforeEach(func() {
			// Needed because the command-table is a singleton
			// and the absence of -v relies on the default value of
			// common.Commands.VerboseOrVersion to be false
			common.Commands.VerboseOrVersion = false
			var err error

			parser, err = command_parser.NewCommandParser()
			Expect(err).ToNot(HaveOccurred())
		})

		It("sets the verbose/version flag", func() {
			status := parser.ParseCommandFromArgs(pluginUI, []string{"-v", "help"})
			Expect(status).To(Equal(0))
			Expect(parser.Config.Flags).To(Equal(configv3.FlagOverride{Verbose: true}))
		})

		It("sets the verbose/version flag after the command-name", func() {
			status := parser.ParseCommandFromArgs(pluginUI, []string{"help", "-v"})
			Expect(status).To(Equal(0))
			Expect(parser.Config.Flags).To(Equal(configv3.FlagOverride{Verbose: true}))
		})

		It("doesn't turn verbose on by default", func() {
			status := parser.ParseCommandFromArgs(pluginUI, []string{"help"})
			Expect(status).To(Equal(0))
			Expect(parser.Config.Flags).To(Equal(configv3.FlagOverride{Verbose: false}))
		})

	})
})
